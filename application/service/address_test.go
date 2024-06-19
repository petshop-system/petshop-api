package service

import (
	"context"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"github.com/petshop-system/petshop-api/adapter/input/http/handler"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/configuration/environment"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
	"time"
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
			Name: "success to save an address",
			Address: domain.AddressDomain{
				Street: "Rua Fulaninho da Silva",
				Number: "123",
			},
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					return domain.AddressDomain{
						ID:     1,
						Street: "Rua Fulaninho da Silva",
						Number: "123",
					}, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.AddressDomain{
				ID:     1,
				Street: "Rua Fulaninho da Silva",
				Number: "123",
			},
			ExpectedError: nil,
		},
		{
			Name: "error to save an address: street is required",
			Address: domain.AddressDomain{
				Street: "",
				Number: "123",
			},
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					return domain.AddressDomain{}, errors.New(handler.StreetIsRequired)
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.AddressDomain{},
			ExpectedError:  errors.New(handler.StreetIsRequired),
		},
		{
			Name: "error to save an address: number is required",
			Address: domain.AddressDomain{
				Street: "Rua Fulaninho da Silva",
				Number: "",
			},
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					return domain.AddressDomain{}, errors.New(handler.NumberIsRequired)
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.AddressDomain{},
			ExpectedError:  errors.New(handler.NumberIsRequired),
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
			Name: "success to get an address by id",
			Address: domain.AddressDomain{
				Street: "Rua Fulaninho da Silva",
				Number: "123",
			},
			AddressDomainDataBaseRepository: output.AddressDomainDataBaseRepositoryMock{
				GetByIDMock: func(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
					return domain.AddressDomain{
						ID:     1,
						Street: "Rua Fulaninho da Silva",
						Number: "123",
					}, true, nil
				},
			},
			AddressDomainCacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.AddressDomain{
				ID:     1,
				Street: "Rua Fulaninho da Silva",
				Number: "123",
			},
			ExpectedExists: true,
			ExpectedError:  nil,
		},
		{
			Name: "address not found",
			Address: domain.AddressDomain{
				Street: "Rua Fulaninho da Silva",
				Number: "123",
			},
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
