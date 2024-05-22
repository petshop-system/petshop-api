package service

import (
	"encoding/json"
	"github.com/petshop-system/petshop-api/application/domain"
	"go.uber.org/zap"
)

type ScheduleService struct {
	LoggerSugar *zap.SugaredLogger
}

func (ss ScheduleService) CreateFromMessage(contextControl domain.ContextControl, scheduleMessage domain.ScheduleMessage) error {

	msg, _ := json.Marshal(scheduleMessage)

	ss.LoggerSugar.Infow("######## Conseguiu buscar", "message", string(msg))

	return nil

}
