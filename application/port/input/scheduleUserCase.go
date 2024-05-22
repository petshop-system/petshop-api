package input

import "github.com/petshop-system/petshop-api/application/domain"

type IScheduleService interface {
	CreateFromMessage(contextControl domain.ContextControl, message domain.ScheduleMessage) error
}
