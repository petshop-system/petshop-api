package service

import "github.com/petshop-system/petshop-api/application/domain"

type PhoneMock struct {
	CreateMock  func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error)
	GetByIDMock func(ID int64) (domain.PhoneDomain, error)
}

func (c PhoneMock) Create(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
	if c.CreateMock != nil {
		return c.CreateMock(contextControl, phone)
	}
	return domain.PhoneDomain{}, nil
}

func (c PhoneMock) GetByID(ID int64) (domain.PhoneDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(ID)
	}
	return domain.PhoneDomain{}, nil
}
