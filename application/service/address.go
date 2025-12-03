package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"go.uber.org/zap"
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
	AddressErrorToSaveInCache    = "error to save an address in cache"
	AddressErrorToGetByIDInCache = "error to getting an address in cache"
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

	hash, err := json.Marshal(save)
	if err != nil {
		service.LoggerSugar.Warnw("failed to marshal address for cache", "address_id", save.ID, "error", err)
	}

	if err = service.AddressDomainCacheRepository.Set(contextControl,
		service.getCacheKey(AddressCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), AddressCacheTTL); err != nil {
		service.LoggerSugar.Infow(AddressErrorToSaveInCache, "address_id", save.ID, "error", err)
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

	hash, err := json.Marshal(address)
	if err != nil {
		service.LoggerSugar.Warnw("failed to marshal address for cache", "address_id", address.ID, "error", err)
	}

	if err = service.AddressDomainCacheRepository.Set(contextControl,
		service.getCacheKey(AddressCacheKeyTypeID, strconv.FormatInt(address.ID, 10)),
		string(hash), AddressCacheTTL); err != nil {
		service.LoggerSugar.Infow(AddressErrorToGetByIDInCache, "address_id", address.ID, "error", err)
	}

	return address, exists, nil
}

// ValidateAddress checks that all required address fields are present and non-empty.
func (service AddressService) ValidateAddress(address domain.AddressDomain) error {
	var errs []error

	if len(strings.TrimSpace(address.Street)) == 0 {
		errs = append(errs, errors.New(StreetIsRequired))
	}
	if len(strings.TrimSpace(address.Number)) == 0 {
		errs = append(errs, errors.New(NumberIsRequired))
	}
	if len(strings.TrimSpace(address.Neighborhood)) == 0 {
		errs = append(errs, errors.New(NeighborhoodIsRequired))
	}
	if len(strings.TrimSpace(address.ZipCode)) == 0 {
		errs = append(errs, errors.New(ZipCodeIsRequired))
	}
	if len(strings.TrimSpace(address.City)) == 0 {
		errs = append(errs, errors.New(CityIsRequired))
	}

	trimmedState := strings.TrimSpace(address.State)
	if len(trimmedState) == 0 {
		errs = append(errs, errors.New(StateIsRequired))
	} else if len(trimmedState) != 2 {
		errs = append(errs, errors.New("state must be exactly 2 characters"))
	}

	if len(strings.TrimSpace(address.Country)) == 0 {
		errs = append(errs, errors.New(CountryIsRequired))
	}

	if len(errs) == 0 {
		return nil
	}

	errorMessages := make([]string, len(errs))
	for i, err := range errs {
		errorMessages[i] = err.Error()
	}
	return fmt.Errorf("error to validate address: %s", strings.Join(errorMessages, ", "))
}
