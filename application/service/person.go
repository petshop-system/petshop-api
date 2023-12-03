package service

import (
	"encoding/json"
	"fmt"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type PersonService struct {
	LoggerSugar                    *zap.SugaredLogger
	PersonDomainDataBaseRepository output.IPersonDomainDataBaseRepository
	PersonDomainCacheRepository    output.IPersonDomainCacheRepository
}

var PersonCacheTTL = 10 * time.Minute

const (
	PersonCacheKeyTypeID = "ID"
)

const (
	PersonErrorToSaveInCache = "error to save person in cache."
)

func (service PersonService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service PersonService) Create(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error) {

	save, err := service.PersonDomainDataBaseRepository.Save(contextControl, person)
	if err != nil {
		return domain.PersonDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.PersonDomainCacheRepository.Set(contextControl,
		service.getCacheKey(PersonCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), PersonCacheTTL); err != nil {
		service.LoggerSugar.Infow(PersonErrorToSaveInCache, "person_id", save.ID)
	}

	return save, nil
}
