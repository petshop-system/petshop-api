package service

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/utils"
	"github.com/petshop-system/petshop-api/configuration/environment"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
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
		{
			Name: "success saving phone",
			Phone: domain.PhoneDomain{
				Number:    "99999-9999",
				CodeArea:  "32",
				PhoneType: utils.MobilePhone,
			},
			PhoneDomainDataBaseRepository: nil,
			PhoneDomainCacheRepository:    nil,
			ExpectedResult:                domain.PhoneDomain{},
			ExpectedError:                 nil,
		},
	}
}
