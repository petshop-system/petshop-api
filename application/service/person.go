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
	TypePersonLegal      = "legal"
	TypePersonIndividual = "individual"
)

const (
	PersonErrorToSaveInCache    = "error to save person in cache"
	PersonErrorToGetByIDInCache = "error to save person in cache"
	InvalidTypeOfDocument       = "invalid type of person"
)

func (service *PersonService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service *PersonService) Create(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error) {

	err := service.Validate(person)
	if err != nil {
		return domain.PersonDomain{}, err
	}

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

func (service *PersonService) GetByID(contextControl domain.ContextControl, ID int64) (domain.PersonDomain, bool, error) {
	person, exists, err := service.PersonDomainDataBaseRepository.GetByID(contextControl, ID)
	if err != nil {
		return domain.PersonDomain{}, exists, err
	}

	if !exists {
		return domain.PersonDomain{}, exists, nil
	}
	hash, _ := json.Marshal(person)
	if err = service.PersonDomainCacheRepository.Set(contextControl,
		service.getCacheKey(AddressCacheKeyTypeID, strconv.FormatInt(person.ID, 10)),
		string(hash), AddressCacheTTL); err != nil {
		service.LoggerSugar.Infow(PersonErrorToGetByIDInCache, "address_id", person.ID)
	}
	return person, exists, nil
}

func (service *PersonService) Validate(person domain.PersonDomain) error {
	switch person.Person_type {
	case TypePersonLegal:
		if err := utils.ValidateCnpj(person.Document); err != nil {
			return err
		}
	case TypePersonIndividual:
		if err := utils.ValidateCpf(person.Document); err != nil {
			return err
		}
	default:
		return fmt.Errorf(InvalidTypeOfDocument)
	}
	return nil
}
