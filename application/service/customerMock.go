package service

import "github.com/petshop-system/petshop-api/application/domain"

type CustomerMock struct {
	CreateMock  func(contextControl domain.ContextControl, customer domain.CustomerDomain) (domain.CustomerDomain, error)
	GetByIDMock func(ID int64) (domain.CustomerDomain, error)
}

func (c CustomerMock) Create(contextControl domain.ContextControl, customer domain.CustomerDomain) (domain.CustomerDomain, error) {
	if c.CreateMock != nil {
		return c.CreateMock(contextControl, customer)
	}
	return domain.CustomerDomain{}, nil
}

func (c CustomerMock) GetByID(ID int64) (domain.CustomerDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(ID)
	}
	return domain.CustomerDomain{}, nil
}
