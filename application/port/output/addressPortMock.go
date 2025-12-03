package output

import (
	"time"

	"github.com/petshop-system/petshop-api/application/domain"
)

type AddressDomainDataBaseRepositoryMock struct {
	SaveMock    func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error)
	GetByIDMock func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error)
}

type AddressDomainCacheRepositoryMock struct {
	SetMock    func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	GetMock    func(contextControl domain.ContextControl, key string) (string, error)
	DeleteMock func(contextControl domain.ContextControl, key string) error
}

func (c AddressDomainDataBaseRepositoryMock) Save(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
	if c.SaveMock != nil {
		return c.SaveMock(contextControl, address)
	}
	return domain.AddressDomain{}, nil
}

func (c AddressDomainDataBaseRepositoryMock) GetByID(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(contextControl, ID)
	}
	return domain.AddressDomain{}, false, nil
}

func (c AddressDomainCacheRepositoryMock) Delete(contextControl domain.ContextControl, key string) error {
	if c.DeleteMock != nil {
		return c.DeleteMock(contextControl, key)
	}
	return nil
}

func (c AddressDomainCacheRepositoryMock) Get(contextControl domain.ContextControl, key string) (string, error) {
	if c.GetMock != nil {
		return c.GetMock(contextControl, key)
	}
	return "", nil
}

func (c AddressDomainCacheRepositoryMock) Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
	if c.SetMock != nil {
		return c.SetMock(contextControl, key, hash, expirationTime)
	}
	return nil
}
