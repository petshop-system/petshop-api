package service

import (
	"context"
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
	defer logger.Sync() // flushes buffer, if any
	loggerSugar = logger.Sugar()
	loggerSugar.Infow("testing phone services")
}

func TestPhoneService_Create(t *testing.T) {

	tests := []struct {
		Name                          string
		Phone                         domain.PhoneDomain
		PhoneDomainDataBaseRepository output.IPhoneDomainDataBaseRepository
		PhoneDomainCacheRepository    output.IPhoneDomainCacheRepository
		ExpectedResult                domain.PhoneDomain
		ExpectedError                 error
	}{
		//Mobile Phone
		{
			Name: "success saving mobile phone",
			Phone: domain.PhoneDomain{
				Number:    "99999-9999",
				CodeArea:  "32",
				PhoneType: MobilePhone,
			},
			PhoneDomainDataBaseRepository: output.PhoneDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
					return domain.PhoneDomain{
						ID:        1,
						Number:    "999999999",
						CodeArea:  "32",
						PhoneType: MobilePhone,
					}, nil
				},
			},
			PhoneDomainCacheRepository: output.PhoneDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.PhoneDomain{
				ID:        1,
				Number:    "999999999",
				CodeArea:  "32",
				PhoneType: MobilePhone,
			},
			ExpectedError: nil,
		},
		{
			Name: "error saving mobile phone",
			Phone: domain.PhoneDomain{
				Number:    "33333399999-9999",
				CodeArea:  "32",
				PhoneType: MobilePhone,
			},
			PhoneDomainDataBaseRepository: output.PhoneDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
					return domain.PhoneDomain{}, fmt.Errorf(database.PhoneSaveError)
				},
			},
			PhoneDomainCacheRepository: output.PhoneDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return fmt.Errorf(PhoneErrorToSaveInCache)
				},
			},
			ExpectedResult: domain.PhoneDomain{},
			ExpectedError:  fmt.Errorf(ErrorInvalidMobilePhoneLength),
		},

		//Landline Phone
		{
			Name: "success saving landline phone",
			Phone: domain.PhoneDomain{
				Number:    "3212-2222",
				CodeArea:  "32",
				PhoneType: LandLinePhone,
			},
			PhoneDomainDataBaseRepository: output.PhoneDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
					return domain.PhoneDomain{
						ID:        1,
						Number:    "32122222",
						CodeArea:  "32",
						PhoneType: LandLinePhone,
					}, nil
				},
			},
			PhoneDomainCacheRepository: output.PhoneDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.PhoneDomain{
				ID:        1,
				Number:    "32122222",
				CodeArea:  "32",
				PhoneType: LandLinePhone,
			},
			ExpectedError: nil,
		},
		{
			Name: "error saving landline phone",
			Phone: domain.PhoneDomain{
				Number:    "00003211-0544",
				CodeArea:  "32",
				PhoneType: LandLinePhone,
			},
			PhoneDomainDataBaseRepository: output.PhoneDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
					return domain.PhoneDomain{}, fmt.Errorf(database.PhoneSaveError)
				},
			},
			PhoneDomainCacheRepository: output.PhoneDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return fmt.Errorf(PhoneErrorToSaveInCache)
				},
			},
			ExpectedResult: domain.PhoneDomain{},
			ExpectedError:  fmt.Errorf(ErrorInvalidLandLinePhoneLength),
		},

		//Code Area
		{
			Name: "success saving code area phone",
			Phone: domain.PhoneDomain{
				Number:    "99999-9999",
				CodeArea:  "96",
				PhoneType: MobilePhone,
			},
			PhoneDomainDataBaseRepository: output.PhoneDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
					return domain.PhoneDomain{
						ID:        1,
						Number:    "999999999",
						CodeArea:  "96",
						PhoneType: MobilePhone,
					}, nil
				},
			},
			PhoneDomainCacheRepository: output.PhoneDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			ExpectedResult: domain.PhoneDomain{
				ID:        1,
				Number:    "999999999",
				CodeArea:  "96",
				PhoneType: MobilePhone,
			},
			ExpectedError: nil,
		},
		{
			Name: "error saving code area phone",
			Phone: domain.PhoneDomain{
				Number:    "99999-9999",
				CodeArea:  "39", //Invalid DDD
				PhoneType: MobilePhone,
			},
			PhoneDomainDataBaseRepository: output.PhoneDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
					return domain.PhoneDomain{}, fmt.Errorf(database.PhoneSaveError)
				},
			},
			PhoneDomainCacheRepository: output.PhoneDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return fmt.Errorf(PhoneErrorToSaveInCache)
				},
			},
			ExpectedResult: domain.PhoneDomain{},
			ExpectedError:  fmt.Errorf(utils.ErrorAreaCodeVerification),
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			phoneService := PhoneService{
				LoggerSugar:                   loggerSugar,
				PhoneDomainCacheRepository:    test.PhoneDomainCacheRepository,
				PhoneDomainDataBaseRepository: test.PhoneDomainDataBaseRepository,
			}

			contextControl := domain.ContextControl{
				Context: context.Background(),
			}

			phone, err := phoneService.Create(contextControl, test.Phone)
			assert.Equal(t, test.ExpectedResult, phone)
			assert.Equal(t, test.ExpectedError, err)

		})
	}
}
