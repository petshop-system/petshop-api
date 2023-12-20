package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	adpterHttpInput "github.com/petshop-system/petshop-api/adapter/input/http"
	"github.com/petshop-system/petshop-api/adapter/input/http/handler"
	"github.com/petshop-system/petshop-api/adapter/output/cache"
	"github.com/petshop-system/petshop-api/adapter/output/database"
	"github.com/petshop-system/petshop-api/application/service"
	"github.com/petshop-system/petshop-api/configuration/environment"
	"github.com/petshop-system/petshop-api/configuration/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"

	"github.com/go-chi/chi"
	"os"
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
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // flushes buffer, if any
	loggerSugar = logger.Sugar()

}

func main() {

	redisCache := cache.NewRedis(loggerSugar)

	postgresConnectionDB := repository.NewPostgresDB(environment.Setting.Postgres.DBUser, environment.Setting.Postgres.DBPassword,
		environment.Setting.Postgres.DBName, environment.Setting.Postgres.DBHost, environment.Setting.Postgres.DBPort, loggerSugar)

	customerPostgresDB := database.NewCustomerPostgresDB(postgresConnectionDB, loggerSugar)
	addressPostgresDB := database.NewAddressPostgresDB(postgresConnectionDB, loggerSugar)
	personPostgresDB := database.NewPersonPostgresDB(postgresConnectionDB, loggerSugar)

	genericHandler := &handler.Generic{
		LoggerSugar: loggerSugar,
	}

	customerService := service.CustomerService{
		LoggerSugar:                      loggerSugar,
		CustomerDomainDataBaseRepository: &customerPostgresDB,
		CustomerDomainCacheRepository:    &redisCache,
	}

	customerHandler := &handler.Customer{
		CustomerService: customerService,
		LoggerSugar:     loggerSugar,
	}

	addressService := service.AddressService{
		LoggerSugar:                     loggerSugar,
		AddressDomainDataBaseRepository: &addressPostgresDB,
		AddressDomainCacheRepository:    &redisCache,
	}

	addressHandler := &handler.Address{
		AddressService: addressService,
		LoggerSugar:    loggerSugar,
	}

	personService := &service.PersonService{
		LoggerSugar:                    loggerSugar,
		PersonDomainDataBaseRepository: &personPostgresDB,
		PersonDomainCacheRepository:    &redisCache,
	}

	personHandler := &handler.Person{
		PersonService: personService,
		LoggerSugar:   loggerSugar,
	}

	contextPath := environment.Setting.Server.Context
	newRouter := adpterHttpInput.GetNewRouter(loggerSugar)
	newRouter.GetChiRouter().Route(fmt.Sprintf("/%s", contextPath), func(r chi.Router) {
		r.NotFound(genericHandler.NotFound)
		r.Group(newRouter.AddGroupHandlerHealthCheck(genericHandler))
		r.Group(newRouter.AddGroupHandlerCustomer(customerHandler))
		r.Group(newRouter.AddGroupHandlerAddress(addressHandler))
		r.Group(newRouter.AddGroupHandlerPerson(personHandler))
	})

	serverHttp := &http.Server{
		Addr:           fmt.Sprintf(":%s", environment.Setting.Server.Port),
		Handler:        newRouter.GetChiRouter(),
		ReadTimeout:    environment.Setting.Server.ReadTimeout,
		WriteTimeout:   environment.Setting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	loggerSugar.Infow("server started", "port", serverHttp.Addr,
		"contextPath", contextPath)

	if err := serverHttp.ListenAndServe(); err != nil {
		loggerSugar.Errorw("error to listen and starts server", "port", serverHttp.Addr,
			"contextPath", contextPath, "err", err.Error())
		panic(err.Error())
	}

}
