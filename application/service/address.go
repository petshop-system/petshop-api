package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"go.uber.org/zap"
	"strconv"
	"strings"
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
	AddressErrorToSaveInCache    = "error to save address in cache"
	AddressErrorToGetByIDInCache = "error to save and address in cache"
)

const (
	StreetIsRequired       = "street is required"
	NumberIsRequired       = "number is required"
	NeighborhoodIsRequired = "neighborhood is required"
	ZipCodeIsRequired      = "zip code is required"
	CityIsRequired         = "city is required"
	StateIsRequired        = "state is required"
	CountryIsRequired      = "country is required"
)

func (service AddressService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service AddressService) Create(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {

	if err := service.ValidateAddress(address); err != nil {
		return domain.AddressDomain{}, err
	}

	save, err := service.AddressDomainDataBaseRepository.Save(contextControl, address)
	if err != nil {
		return domain.AddressDomain{}, err
	}

	hash, _ := json.Marshal(save)

	if err = service.AddressDomainCacheRepository.Set(contextControl,
		service.getCacheKey(AddressCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), AddressCacheTTL); err != nil {
		service.LoggerSugar.Infow(AddressErrorToSaveInCache, "address_id", save.ID)
		return domain.AddressDomain{}, err
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
		return domain.AddressDomain{}, exists, err
	}

	return address, exists, nil
}

func (service AddressService) ValidateAddress(address domain.AddressDomain) error {
	if len(strings.TrimSpace(address.Street)) == 0 {
		return errors.New(StreetIsRequired)
	}
	if len(strings.TrimSpace(address.Number)) == 0 {
		return errors.New(NumberIsRequired)
	}
	if len(strings.TrimSpace(address.Neighborhood)) == 0 {
		return errors.New(NeighborhoodIsRequired)
	}
	if len(strings.TrimSpace(address.ZipCode)) == 0 {
		return errors.New(ZipCodeIsRequired)
	}
	if len(strings.TrimSpace(address.City)) == 0 {
		return errors.New(CityIsRequired)
	}
	if len(strings.TrimSpace(address.State)) == 0 {
		return errors.New(StateIsRequired)
	}
	if len(strings.TrimSpace(address.Country)) == 0 {
		return errors.New(CountryIsRequired)
	}
	return nil
}
