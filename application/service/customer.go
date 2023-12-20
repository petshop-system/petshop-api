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

type CustomerService struct {
	LoggerSugar                      *zap.SugaredLogger
	CustomerDomainDataBaseRepository output.ICustomerDomainDataBaseRepository
	CustomerDomainCacheRepository    output.ICustomerDomainCacheRepository
}

var ClientCacheTTL = 10 * time.Minute

const (
	CustomerCacheKeyTypeID = "ID"
)

const (
	CustomerErrorToSaveInCache = "error to save customer in cache."
)

func (service CustomerService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service CustomerService) Create(contextControl domain.ContextControl, customer domain.CustomerDomain) (domain.CustomerDomain, error) {

	customer.DateCreated = time.Now()
	save, err := service.CustomerDomainDataBaseRepository.Save(contextControl, customer)
	if err != nil {
		return domain.CustomerDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.CustomerDomainCacheRepository.Set(contextControl,
		service.getCacheKey(CustomerCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), ClientCacheTTL); err != nil {
		service.LoggerSugar.Infow(CustomerErrorToSaveInCache, "customer_id", save.ID)
	}

	return save, nil
}
