package service

import (
	"context"
	"github.com/kelseyhightower/envconfig"
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

var loggerSugar *zap.SugaredLogger

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
	loggerSugar.Infow("testing client services")

}

func TestCustomerService_Create(t *testing.T) {

	tests := []struct {
		Name                             string
		Customer                         domain.CustomerDomain
		CustomerDomainDataBaseRepository output.ICustomerDomainDataBaseRepository
		CustomerDomainCacheRepository    output.ICustomerDomainCacheRepository
		ExpectedResult                   domain.CustomerDomain
		ExpectedError                    error
	}{
		{
			Name: "success saving customer",
			Customer: domain.CustomerDomain{
				Name: "Fulano",
			},
			CustomerDomainDataBaseRepository: output.CustomerDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, customer domain.CustomerDomain) (domain.CustomerDomain, error) {
					return domain.CustomerDomain{
						ID:   1,
						Name: "Fulano",
					}, nil
				},
			},
			CustomerDomainCacheRepository: output.CustomerDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.CustomerDomain{
				ID:   1,
				Name: "Fulano",
			},
			ExpectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			customerService := CustomerService{
				LoggerSugar:                      loggerSugar,
				CustomerDomainCacheRepository:    test.CustomerDomainCacheRepository,
				CustomerDomainDataBaseRepository: test.CustomerDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			customer, err := customerService.Create(contextControl, test.Customer)
			assert.Equal(t, test.ExpectedResult, customer)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}
