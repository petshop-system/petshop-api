package output

import (
	"github.com/petshop-system/petshop-api/application/domain"
	"time"
)

type PhoneDomainDataBaseRepositoryMock struct {
	SaveMock    func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error)
	GetByIDMock func(contextControl domain.ContextControl, ID int64) (domain.PhoneDomain, bool, error)
}

type PhoneDomainCacheRepositoryMock struct {
	SetMock    func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	GetMock    func(contextControl domain.ContextControl, key string) (string, error)
	DeleteMock func(contextControl domain.ContextControl, key string) error
}

func (c PhoneDomainDataBaseRepositoryMock) Save(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
	if c.SaveMock != nil {
		return c.SaveMock(contextControl, phone)
	}
	return domain.PhoneDomain{}, nil
}

func (c PhoneDomainDataBaseRepositoryMock) GetByID(contextControl domain.ContextControl, ID int64) (domain.PhoneDomain, bool, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(contextControl, ID)
	}
	return domain.PhoneDomain{}, false, nil
}

func (c PhoneDomainCacheRepositoryMock) Delete(contextControl domain.ContextControl, key string) error {
	if c.DeleteMock != nil {
		return c.DeleteMock(contextControl, key)
	}
	return nil
}

func (c PhoneDomainCacheRepositoryMock) Get(contextControl domain.ContextControl, key string) (string, error) {
	if c.GetMock != nil {
		return c.GetMock(contextControl, key)
	}
	return "", nil
}

func (c PhoneDomainCacheRepositoryMock) Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
	if c.SetMock != nil {
		return c.SetMock(contextControl, key, hash, expirationTime)
	}
	return nil
}
