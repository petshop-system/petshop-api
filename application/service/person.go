package service

import (
	"encoding/json"
	"fmt"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/utils"
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
	TypePersonLegal      = "pessoa_juridica"
	TypePersonIndividual = "pessoa_fisica"
)

const (
	PersonErrorToSaveInCache = "error to save person in cache"
)

func (service PersonService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service PersonService) Create(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error) {

	if person.Person_type == TypePersonIndividual {
		if err := utils.ValidateCpf(person.Document); err != nil {
			return domain.PersonDomain{}, err
		}
	}

	/*if person.Person_type == TypePersonLegal {
		if err := utils.ValidateCnpj(person.Document); err != nil {
			return domain.PersonDomain{}, err
		}
	}*/

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
