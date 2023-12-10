package service

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/petshop-system/petshop-api/adapter/output/database"
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

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // flushes buffer, if any
	loggerSugar = logger.Sugar()
	loggerSugar.Infow("testing person services")

}

func TestPersonService_Create(t *testing.T) {

	tests := []struct {
		Name                           string
		Person                         domain.PersonDomain
		PersonDomainDataBaseRepository output.IPersonDomainDataBaseRepository
		PersonDomainCacheRepository    output.IPersonDomainCacheRepository
		ExpectedResult                 domain.PersonDomain
		ExpectedError                  error
	}{
		{
			Name: "success saving person",
			Person: domain.PersonDomain{
				Document:    "076.164.346-06",
				Person_type: TypePersonIndividual,
			},
			PersonDomainDataBaseRepository: output.PersonDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error) {
					return domain.PersonDomain{
						ID:          1,
						Document:    "076.164.346-06",
						Person_type: TypePersonIndividual,
					}, nil
				},
			},
			PersonDomainCacheRepository: output.PersonDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.PersonDomain{
				ID:          1,
				Document:    "076.164.346-06",
				Person_type: TypePersonIndividual,
			},
			ExpectedError: nil,
		},
		{
			Name: "error saving person",
			Person: domain.PersonDomain{
				Document:    "076.164.346-06",
				Person_type: TypePersonIndividual,
			},
			PersonDomainDataBaseRepository: output.PersonDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, person domain.PersonDomain) (domain.PersonDomain, error) {
					return domain.PersonDomain{}, fmt.Errorf(database.PersonSaveDBError)
				},
			},
			ExpectedResult: domain.PersonDomain{},
			ExpectedError:  fmt.Errorf(database.PersonSaveDBError),
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			personService := PersonService{
				LoggerSugar:                    loggerSugar,
				PersonDomainCacheRepository:    test.PersonDomainCacheRepository,
				PersonDomainDataBaseRepository: test.PersonDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			person, err := personService.Create(contextControl, test.Person)
			assert.Equal(t, test.ExpectedResult, person)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}
