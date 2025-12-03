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
	//logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // flushes buffer, if any
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
			Name: "Test Successful - Address correctly saved to the database",
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
			Name: "Test Failure - Error while validating address",
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
			Name: "Test Failure - Error while saving address to the database",
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
			Name: "Test Failure - Error while saving address to the cache",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
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
			ExpectedResult: domain.AddressDomain{},
			ExpectedError:  errors.New(AddressErrorToSaveInCache),
		},
		{
			Name: "Test Failure - Multiple validation errors",
			Address: func() domain.AddressDomain {
				return domain.AddressDomain{} // Todos os campos vazios
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
			Name: "Test Successful - Getting an address by id",
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
			Name: "Test Failure - Error to get an address by id",
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
			Name: "Test Failure - Error returned from repository",
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
			Name: "Test Failure - Error returned from cache repository",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
					return utils.GetMockAddress(), true, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return errors.New(CacheError)
				},
			},
			ExpectedResult: domain.AddressDomain{},
			ExpectedExists: true,
			ExpectedError:  errors.New(CacheError),
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
			Name: "Test Successful - Validating address",
			Address: func() domain.AddressDomain {
				return utils.GetMockAddress()
			}(),
			ExpectedError: nil,
		},
		{
			Name: "Test Failure - Error to validate address: street is required",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Street = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", StreetIsRequired),
		},

		{
			Name: "Test Failure - Error to validate address: number is required",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Number = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", NumberIsRequired),
		},
		{
			Name: "Test Failure - Error to validate address: neighborhood is required",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.Neighborhood = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", NeighborhoodIsRequired),
		},
		{
			Name: "Test Failure - Error to validate address: zip code is required",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.ZipCode = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", ZipCodeIsRequired),
		},
		{
			Name: "Test Failure - Error to validate address: city is required",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.City = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", CityIsRequired),
		},
		{
			Name: "Test Failure - Error to validate address: state is required",
			Address: func() domain.AddressDomain {
				address := utils.GetMockAddress()
				address.State = ""
				return address
			}(),
			ExpectedError: fmt.Errorf("error to validate address: %s", StateIsRequired),
		},
		{
			Name: "Test Failure - Error to validate address: country is required",
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
