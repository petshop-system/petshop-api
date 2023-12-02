package service

import (
	"encoding/json"
	"fmt"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type ClienteService struct {
	LoggerSugar                      *zap.SugaredLogger
	CustomerDomainDataBaseRepository output.ICustomerDomainDataBaseRepository
	CustomerDomainCacheRepository    output.ICustomerDomainCacheRepository
}

var ClientCacheTTL = 10 * time.Minute

const (
	ClienteCacheKeyTypeID = "ID"
)

const (
	ClienteErrorToSaveInCache = "error to save cliente in cache."
)

func (service ClienteService) getCacheKey(cacheKeyType string, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service ClienteService) Create(contextControl domain.ContextControl, cliente domain.ClienteDomain) (domain.ClienteDomain, error) {

	cliente.DataCadastro = time.Now()
	save, err := service.CustomerDomainDataBaseRepository.Save(contextControl, cliente)
	if err != nil {
		return domain.ClienteDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.CustomerDomainCacheRepository.Set(contextControl,
		service.getCacheKey(ClienteCacheKeyTypeID, strconv.FormatInt(save.ID, 10)),
		string(hash), ClientCacheTTL); err != nil {
		service.LoggerSugar.Infow(ClienteErrorToSaveInCache, "cliente_id", save.ID)
	}

	return save, nil
}
