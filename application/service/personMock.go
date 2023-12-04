package service

import "github.com/petshop-system/petshop-api/application/domain"

type PersonMock struct {
	CreateMock  func(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error)
	GetByIDMock func(ID int64) (domain.PersonDomain, error)
}

func (c PersonMock) Create(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error) {
	if c.CreateMock != nil {
		return c.CreateMock(contextControl, person)
	}
	return domain.PersonDomain{}, nil
}

func (c PersonMock) GetByID(ID int64) (domain.PersonDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(ID)
	}
	return domain.PersonDomain{}, nil
}
