package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/input"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
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
	KafkaClient     *kgo.Client
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

	seeds := []string{bootstrapServer}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
	)

	if err != nil {
		loggerSugar.Errorw(ScheduleKafkaErrorToStartConsumer, "error", err.Error())
		panic(err)
	}

	scheduleKafkaConsumer := ScheduleKafkaConsumer{
		ScheduleService: scheduleService,
		LoggerSugar:     loggerSugar,
		KafkaClient:     kafkaClient,
	}

	return scheduleKafkaConsumer
}

func (schedule *ScheduleKafkaConsumer) ConsumerMessages() {

	go func() {

		for {

			ctx := context.Background()
			fetches := schedule.KafkaClient.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				// All errors are retried internally when fetching, but non-retriable errors are
				// returned from polls so that users can notice and take action.
				//panic(fmt.Sprint(errs))
				schedule.LoggerSugar.Errorw(ScheduleKafkaConsumerErrorToReadMessage, "error", fmt.Sprint(errs))
				continue
			}

			// We can iterate through a record iterator...
			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				//fmt.Println(string(record.Value), "from an iterator!")

				var scheduleMessageKafka ScheduleMessageKafka
				json.NewDecoder(bytes.NewReader(record.Value)).Decode(&scheduleMessageKafka)

				var scheduleMessage domain.ScheduleMessage
				copier.Copy(&scheduleMessage, &scheduleMessageKafka)
				schedule.ScheduleService.CreateFromMessage(domain.ContextControl{
					Context: context.Background(),
				}, scheduleMessage)

				schedule.LoggerSugar.Infow(ScheduleKafkaConsumerSuccessToConsumer,
					"message", string(record.Value))

			}

			// or a callback function.
			//fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			//	for _, record := range p.Records {
			//		fmt.Println(string(record.Value), "from range inside a callback!")
			//	}
			//
			//	// We can even use a second callback!
			//	p.EachRecord(func(record *kgo.Record) {
			//		fmt.Println(string(record.Value), "from a second callback!")
			//	})
			//})
		}
	}()
}
