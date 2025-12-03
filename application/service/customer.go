package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/utils"
	"go.uber.org/zap"
)

type CustomerService struct {
	LoggerSugar                      *zap.SugaredLogger
	CustomerDomainDataBaseRepository output.ICustomerDomainDataBaseRepository
	CustomerDomainCacheRepository    output.ICustomerDomainCacheRepository
}

var ClientCacheTTL = 10 * time.Minute

const (
	CustomerCacheKeyTypeID = "CUSTOMER_ID"
)

const (
	CustomerErrorToSaveInCache    = "error to save customer in cache."
	CustomerErrorToGetByIDInCache = "error to get person in cache"
	InvalidTypeOfDocument         = "invalid type of person"
)

const (
	TypePersonLegal      = "legal"
	TypePersonIndividual = "individual"
)

func (service *CustomerService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service *CustomerService) Create(contextControl domain.ContextControl, customer domain.CustomerDomain) (domain.CustomerDomain, error) {

	err := service.ValidateTypePerson(customer)
	if err != nil {
		return domain.CustomerDomain{}, err
	}

	customer.Document = utils.RemoveNonAlphaNumericCharacters(customer.Document)
	save, err := service.CustomerDomainDataBaseRepository.Save(contextControl, customer)
	if err != nil {
		return domain.CustomerDomain{}, err
	}

	hash, err := json.Marshal(save)
	if err != nil {
		service.LoggerSugar.Warnw("failed to marshal customer for cache", "customer_id", save.ID, "error", err)
	}
	if err = service.CustomerDomainCacheRepository.Set(contextControl,
		service.getCacheKey(CustomerCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), ClientCacheTTL); err != nil {
		service.LoggerSugar.Infow(CustomerErrorToSaveInCache, "customer_id", save.ID)
	}

	return save, nil
}

func (service *CustomerService) ValidateTypePerson(customer domain.CustomerDomain) error { //TODO: Change the method name to ValidatePerson
	switch customer.PersonType {
	case TypePersonLegal:
		if err := utils.ValidateCnpj(customer.Document); err != nil {
			return err
		}
	case TypePersonIndividual:
		if err := utils.ValidateCpf(customer.Document); err != nil {
			return err
		}
	default:
		return fmt.Errorf(InvalidTypeOfDocument)
	}
	return nil
}

func (service *CustomerService) ValidateCreate(customer domain.CustomerDomain) error {

	if err := service.ValidateTypePerson(customer); err != nil {
		return err
	}

	return nil
}
