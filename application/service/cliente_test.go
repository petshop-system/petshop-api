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
	loggerSugar.Infow("testing")

}

func TestCustomerService_Create(t *testing.T) {

	tests := []struct {
		Name                             string
		Cliente                          domain.ClienteDomain
		CustomerDomainDataBaseRepository output.ICustomerDomainDataBaseRepository
		CustomerDomainCacheRepository    output.ICustomerDomainCacheRepository
		ExpectedResult                   domain.ClienteDomain
		ExpectedError                    error
	}{
		{
			Name: "sucesso ao salvar cliente",
			Cliente: domain.ClienteDomain{
				Nome: "Fulano",
			},
			CustomerDomainDataBaseRepository: output.CustomerDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, cliente domain.ClienteDomain) (domain.ClienteDomain, error) {
					return domain.ClienteDomain{
						ID:   1,
						Nome: "Fulano",
					}, nil
				},
			},
			CustomerDomainCacheRepository: output.CustomerDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.ClienteDomain{
				ID:   1,
				Nome: "Fulano",
			},
			ExpectedError: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			clienteService := ClienteService{
				LoggerSugar:                      loggerSugar,
				CustomerDomainCacheRepository:    test.CustomerDomainCacheRepository,
				CustomerDomainDataBaseRepository: test.CustomerDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			cliente, err := clienteService.Create(contextControl, test.Cliente)
			assert.Equal(t, test.ExpectedResult, cliente)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}
