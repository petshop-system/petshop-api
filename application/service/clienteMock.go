package service

import "github.com/petshop-system/petshop-api/application/domain"

type CustomerMock struct {
	CreateMock  func(contextControl domain.ContextControl, cliente domain.ClienteDomain) (domain.ClienteDomain, error)
	GetByIDMock func(ID int64) (domain.ClienteDomain, error)
}

func (c CustomerMock) Create(contextControl domain.ContextControl, cliente domain.ClienteDomain) (domain.ClienteDomain, error) {
	if c.CreateMock != nil {
		return c.CreateMock(contextControl, cliente)
	}
	return domain.ClienteDomain{}, nil
}

func (c CustomerMock) GetByID(ID int64) (domain.ClienteDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(ID)
	}
	return domain.ClienteDomain{}, nil
}
