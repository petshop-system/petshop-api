package output

import (
	"time"

	"github.com/petshop-system/petshop-api/application/domain"
)

type ICustomerDomainDataBaseRepository interface {
	Save(contextControl domain.ContextControl, customer domain.CustomerDomain) (domain.CustomerDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.CustomerDomain, error)
}

type ICustomerDomainCacheRepository interface {
	Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	Get(contextControl domain.ContextControl, key string) (string, error)
	Delete(contextControl domain.ContextControl, key string) error
}
