package output

import (
	"github.com/petshop-system/petshop-api/application/domain"
	"time"
)

type ICustomerDomainDataBaseRepository interface {
	Save(contextControl domain.ContextControl, cliente domain.ClienteDomain) (domain.ClienteDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.ClienteDomain, error)
}

type ICustomerDomainCacheRepository interface {
	Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	Get(contextControl domain.ContextControl, key string) (string, error)
	Delete(contextControl domain.ContextControl, key string) error
}
