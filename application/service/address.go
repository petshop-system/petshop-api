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

type AddressService struct {
	LoggerSugar                     *zap.SugaredLogger
	AddressDomainDataBaseRepository output.IAddressDomainDataBaseRepository
	AddressDomainCacheRepository    output.IAddressDomainCacheRepository
}

var AddressCacheTTL = 10 * time.Minute

const (
	AddressCacheKeyTypeID = "ID"
)

const (
	AddressErrorToSaveInCache    = "error to save address in cache."
	AddressErrorToGetByIDInCache = "error to save and address in cache"
)

func (service AddressService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service AddressService) Create(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {

	save, err := service.AddressDomainDataBaseRepository.Save(contextControl, address)
	if err != nil {
		return domain.AddressDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.AddressDomainCacheRepository.Set(contextControl,
		service.getCacheKey(AddressCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), AddressCacheTTL); err != nil {
		service.LoggerSugar.Infow(AddressErrorToSaveInCache, "address_id", save.ID)
	}

	return save, nil
}

func (service AddressService) GetByID(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
	address, exists, err := service.AddressDomainDataBaseRepository.GetByID(contextControl, ID)
	if err != nil {
		return domain.AddressDomain{}, exists, err
	}

	if !exists {
		return domain.AddressDomain{}, exists, nil
	}
	hash, _ := json.Marshal(address)
	if err = service.AddressDomainCacheRepository.Set(contextControl,
		service.getCacheKey(AddressCacheKeyTypeID, strconv.FormatInt(address.ID, 10)),
		string(hash), AddressCacheTTL); err != nil {
		service.LoggerSugar.Infow(AddressErrorToGetByIDInCache, "address_id", address.ID)
	}

	return address, exists, nil
}
