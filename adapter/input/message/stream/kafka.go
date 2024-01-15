package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/input"
	"go.uber.org/zap"
	"time"
)

const (
	ScheduleKafkaConsumerErrorToReadMessage        = "error to read message from schedule kafka consumer"
	ScheduleKafkaConsumerErrorTimeoutToReadMessage = "timeout error to read message from schedule kafka consumer"
	ScheduleKafkaConsumerSuccessToConsumer         = "success to consumer"
	ScheduleKafkaErrorToStartConsumer              = "error to start consumer from kafka"
)

type ScheduleKafkaConsumer struct {
	LoggerSugar     *zap.SugaredLogger
	ScheduleService input.IScheduleService
	Consumer        *kafka.Consumer
}

type ScheduleMessageKafka struct {
	Booking                    string `json:"booking"`
	PetId                      int    `json:"pet_id"`
	ServiceEmployeeAttentionId int    `json:"service_employee_attention_id"`
}

func NewScheduleKafkaClient(loggerSugar *zap.SugaredLogger,
	scheduleService input.IScheduleService,
	bootstrapServer string,
	groupID string,
	autoOffsetReset string,
	topic string) ScheduleKafkaConsumer {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          groupID,
		"auto.offset.reset": autoOffsetReset,
	})

	if err != nil {
		loggerSugar.Errorw(ScheduleKafkaErrorToStartConsumer, "error", err.Error())
		panic(err)
	}

	consumer.SubscribeTopics([]string{topic}, nil)

	scheduleKafkaConsumer := ScheduleKafkaConsumer{
		ScheduleService: scheduleService,
		LoggerSugar:     loggerSugar,
		Consumer:        consumer,
	}

	return scheduleKafkaConsumer
}

func (schedule *ScheduleKafkaConsumer) ConsumerMessages() {

	go func() {

		defer schedule.Consumer.Close()

		for true {
			message, err := schedule.Consumer.ReadMessage(time.Second * 2)
			if message == nil || len(string(message.Value)) == 0 {
				continue
			}
			schedule.LoggerSugar.With("topic", message.TopicPartition)
			if err != nil {
				schedule.LoggerSugar.Errorw(ScheduleKafkaConsumerErrorToReadMessage, "error", err.Error(),
					"message", string(message.Value))
				continue
			}

			if !err.(kafka.Error).IsTimeout() {
				schedule.LoggerSugar.Errorw(ScheduleKafkaConsumerErrorTimeoutToReadMessage, "error", err.Error(),
					"message", message)
				continue
			}

			var scheduleMessageKafka ScheduleMessageKafka
			json.NewDecoder(bytes.NewReader(message.Value)).Decode(scheduleMessageKafka)
			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			var scheduleMessage domain.ScheduleMessage
			copier.Copy(scheduleMessage, scheduleMessageKafka)
			schedule.ScheduleService.CreateFromMessage(contextControl, domain.ScheduleMessage{})

			schedule.LoggerSugar.Infow(ScheduleKafkaConsumerSuccessToConsumer, "message", string(message.Value))

		}
	}()
}
