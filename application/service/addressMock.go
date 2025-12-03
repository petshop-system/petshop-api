package service

import "github.com/petshop-system/petshop-api/application/domain"

// AddressMock provides a mock implementation of the address service for testing purposes.
// Each method can be overridden by setting the corresponding Mock field to a custom function.
type AddressMock struct {
	CreateMock  func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error)
	GetByIDMock func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error)
}

func (c AddressMock) Create(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
	if c.CreateMock != nil {
		return c.CreateMock(contextControl, address)
	}
	return domain.AddressDomain{}, nil
}

func (c AddressMock) GetByID(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(contextControl, ID)
	}
	return domain.AddressDomain{}, false, nil
}
