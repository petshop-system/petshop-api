package output

import (
	"github.com/petshop-system/petshop-api/application/domain"
	"time"
)

type IPersonDomainDataBaseRepository interface {
	Save(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.PersonDomain, error)
}

type IPersonDomainCacheRepository interface {
	Set(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error
	Get(contextControl domain.ContextControl, key string) (string, error)
	Delete(contextControl domain.ContextControl, key string) error
}
