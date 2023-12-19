package output

import (
	"github.com/petshop-system/petshop-api/application/domain"
	"time"
)

type PersonDomainDataBaseRepositoryMock struct {
	SaveMock    func(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error)
	GetByIDMock func(contextControl domain.ContextControl, ID int64) (domain.PersonDomain, bool, error)
}

type PersonDomainCacheRepositoryMock struct {
	SetMock    func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	GetMock    func(contextControl domain.ContextControl, key string) (string, error)
	DeleteMock func(contextControl domain.ContextControl, key string) error
}

func (c PersonDomainDataBaseRepositoryMock) Save(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error) {
	if c.SaveMock != nil {
		return c.SaveMock(contextControl, person)
	}
	return domain.PersonDomain{}, nil
}

func (c PersonDomainDataBaseRepositoryMock) GetByID(contextControl domain.ContextControl, ID int64) (domain.PersonDomain, bool, error) {
	if c.GetByIDMock != nil {
		return c.GetByIDMock(contextControl, ID)
	}
	return domain.PersonDomain{}, false, nil
}

func (c PersonDomainCacheRepositoryMock) Delete(contextControl domain.ContextControl, key string) error {
	if c.DeleteMock != nil {
		return c.DeleteMock(contextControl, key)
	}
	return nil
}

func (c PersonDomainCacheRepositoryMock) Get(contextControl domain.ContextControl, key string) (string, error) {
	if c.GetMock != nil {
		return c.GetMock(contextControl, key)
	}
	return "", nil
}

func (c PersonDomainCacheRepositoryMock) Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
	if c.SetMock != nil {
		return c.SetMock(contextControl, key, hash, expirationTime)
	}
	return nil
}
