package output

import (
	"github.com/petshop-system/petshop-api/application/domain"
	"time"
)

type CustomerDomainDataBaseRepositoryMock struct {
	SaveMock    func(contextControl domain.ContextControl, cliente domain.ClienteDomain) (domain.ClienteDomain, error)
	GetByIDMock func(contextControl domain.ContextControl, ID int64) (domain.ClienteDomain, error)
}

type CustomerDomainCacheRepositoryMock struct {
	SetMock    func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	GetMock    func(contextControl domain.ContextControl, key string) (string, error)
	DeleteMock func(contextControl domain.ContextControl, key string) error
}

func (c CustomerDomainDataBaseRepositoryMock) Save(contextControl domain.ContextControl, cliente domain.ClienteDomain) (domain.ClienteDomain, error) {
	if c.SaveMock != nil {
		return c.SaveMock(contextControl, cliente)
	}
	return domain.ClienteDomain{}, nil
}

func (c CustomerDomainDataBaseRepositoryMock) GetByID(contextControl domain.ContextControl, ID int64) (domain.ClienteDomain, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(contextControl, ID)
	}
	return domain.ClienteDomain{}, nil
}

func (c CustomerDomainCacheRepositoryMock) Delete(contextControl domain.ContextControl, key string) error {
	if c.DeleteMock != nil {
		return c.DeleteMock(contextControl, key)
	}
	return nil
}

func (c CustomerDomainCacheRepositoryMock) Get(contextControl domain.ContextControl, key string) (string, error) {
	if c.GetMock != nil {
		return c.GetMock(contextControl, key)
	}
	return "", nil
}

func (c CustomerDomainCacheRepositoryMock) Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
	if c.SetMock != nil {
		return c.SetMock(contextControl, key, hash, expirationTime)
	}
	return nil
}
