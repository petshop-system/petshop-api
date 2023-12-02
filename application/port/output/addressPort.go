package output

import (
	"github.com/petshop-system/petshop-api/application/domain"
	"time"
)

type IAddressDomainDataBaseRepository interface {
	Save(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, error)
}

type IAddressDomainCacheRepository interface {
	Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	Get(contextControl domain.ContextControl, key string) (string, error)
	Delete(contextControl domain.ContextControl, key string) error
}
