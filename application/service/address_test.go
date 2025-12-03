package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/petshop-system/petshop-api/adapter/output/database"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/utils"
	"github.com/petshop-system/petshop-api/configuration/environment"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {

	err := envconfig.Process("setting", &environment.Setting)
	if err != nil {
		panic(err.Error())
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	defer func() {
		_ = logger.Sync()
	}()
	loggerSugar = logger.Sugar()
	loggerSugar.Infow("testing address services")

}

func TestAddressService_Create(t *testing.T) {

	tests := []struct {
		Name                            string
		Address                         domain.AddressDomain
		AddressDomainDataBaseRepository output.IAddressDomainDataBaseRepository
		AddressDomainCacheRepository    output.IAddressDomainCacheRepository
		ExpectedResult                  domain.AddressDomain
		ExpectedError                   error
	}{
		{
			Name: "WithValidAddress_SavesSuccessfully",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					address = utils.GetMockAddress()
					address.ID = 1
					return address, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.ID = 1
				return address
			}(),
			ExpectedError: nil,
		},
		{
			Name: "WithInvalidAddress_ValidationFails",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Street = ""
				return address
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{},
			AddressDomainCacheRepository:    output.AddressDomainCacheRepositoryMock{},
			ExpectedResult:                  domain.AddressDomain{},
			ExpectedError:                   fmt.Errorf("error to validate address: %s", StreetIsRequired),
		},
		{
			Name: "WithDatabaseError_ReturnsSaveError",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()

			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					return domain.AddressDomain{}, errors.New(database.AddressSaveDBError)
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.AddressDomain{},
			ExpectedError:  errors.New(database.AddressSaveDBError),
		},
		{
			Name: "WithCacheError_SucceedsWithLog",
			Address: func() domain.AddressDomain {
				addr := utils.GetMockAddress()
				addr.ID = 3
				return addr
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					return address, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return errors.New(AddressErrorToSaveInCache)
				},
			},
			ExpectedResult: func() domain.AddressDomain {
				addr := utils.GetMockAddress()
				addr.ID = 3
				return addr
			}(),
			ExpectedError: nil, // Cache failure is non-fatal
		},
		{
			Name: "WithMultipleValidationErrors_ReturnsAllErrors",
			Address: func() domain.AddressDomain {
				return domain.AddressDomain{}
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s, %s, %s, %s, %s, %s, %s",
				StreetIsRequired, NumberIsRequired, NeighborhoodIsRequired,
				ZipCodeIsRequired, CityIsRequired, StateIsRequired, CountryIsRequired),
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			addressService := AddressService{
				LoggerSugar:                     loggerSugar,
				AddressDomainCacheRepository:    test.AddressDomainCacheRepository,
				AddressDomainDataBaseRepository: test.AddressDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			address, err := addressService.Create(contextControl, test.Address)
			assert.Equal(t, test.ExpectedResult, address)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}

func TestAddressService_GetById(t *testing.T) {

	const CacheError = "cache error"

	tests := []struct {
		Name                            string
		Address                         domain.AddressDomain
		AddressDomainDataBaseRepository output.IAddressDomainDataBaseRepository
		AddressDomainCacheRepository    output.IAddressDomainCacheRepository
		ExpectedResult                  domain.AddressDomain
		ExpectedExists                  bool
		ExpectedError                   error
	}{
		{
			Name: "WithValidID_ReturnsAddress",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
					address := utils.GetMockAddress()
					address.ID = 1
					return address, true, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.ID = 1
				return address
			}(),
			ExpectedExists: true,
			ExpectedError:  nil,
		},
		{
			Name: "WithNonExistentID_ReturnsNotFound",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
					return domain.AddressDomain{}, false, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.AddressDomain{},
			ExpectedExists: false,
			ExpectedError:  nil,
		},
		{
			Name: "WithRepositoryError_ReturnsError",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
					return domain.AddressDomain{}, false, errors.New(database.AddressNotFound)
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.AddressDomain{},
			ExpectedExists: false,
			ExpectedError:  errors.New(database.AddressNotFound),
		},
		{
			Name: "WithCacheError_ReturnsAddressSuccessfully",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
					address := utils.GetMockAddress()
					return address, true, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return errors.New(CacheError)
				},
			},
			ExpectedResult: utils.GetMockAddress(),
			ExpectedExists: true,
			ExpectedError:  nil,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			addressService := AddressService{
				LoggerSugar:                     loggerSugar,
				AddressDomainCacheRepository:    test.AddressDomainCacheRepository,
				AddressDomainDataBaseRepository: test.AddressDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			address, exists, err := addressService.GetByID(contextControl, 1)
			assert.Equal(t, test.ExpectedResult, address)
			assert.Equal(t, test.ExpectedExists, exists)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}

func TestAddressService_ValidateAddress(t *testing.T) {

	tests := []struct {
		Name          string
		Address       domain.AddressDomain
		ExpectedError error
	}{
		{
			Name: "WithValidAddress_ReturnsNoError",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			ExpectedError: nil,
		},
		{
			Name: "WithEmptyStreet_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Street = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", StreetIsRequired),
		},

		{
			Name: "WithEmptyNumber_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Number = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", NumberIsRequired),
		},
		{
			Name: "WithEmptyNeighborhood_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Neighborhood = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", NeighborhoodIsRequired),
		},
		{
			Name: "WithEmptyZipCode_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.ZipCode = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", ZipCodeIsRequired),
		},
		{
			Name: "WithEmptyCity_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.City = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", CityIsRequired),
		},
		{
			Name: "WithEmptyState_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.State = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", StateIsRequired),
		},
		{
			Name: "WithInvalidStateLength_OneCharacter_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.State = "R" // 1 character instead of 2
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: state must be exactly 2 characters"),
		},
		{
			Name: "WithInvalidStateLength_ThreeCharacters_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.State = "RJJ" // 3 characters instead of 2
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: state must be exactly 2 characters"),
		},
		{
			Name: "WithEmptyCountry_ReturnsError",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Country = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", CountryIsRequired),
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			addressService := AddressService{
				AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{},
			}

			err := addressService.ValidateAddress(test.Address)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}
