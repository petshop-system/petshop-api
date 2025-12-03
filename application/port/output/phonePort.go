package output

import (
	"time"

	"github.com/petshop-system/petshop-api/application/domain"
)

type IPhoneDomainDataBaseRepository interface {
	Save(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.PhoneDomain, bool, error)
}

type IPhoneDomainCacheRepository interface {
	Set(contextControl domain.ContextControl, key, hash string, expirationTime time.Duration) error
	Get(contextControl domain.ContextControl, key string) (string, error)
	Delete(contextControl domain.ContextControl, key string) error
}
