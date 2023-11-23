package service

import (
	"github.com/petshop-system/petshop-api/application/domain"
	"go.uber.org/zap"
)

type CustomerService struct {
	LoggerSugar *zap.SugaredLogger
}

func (service CustomerService) Create(contextControl domain.ContextControl, customer domain.Customer) (domain.Customer, error) {

	return domain.Customer{}, nil
}
