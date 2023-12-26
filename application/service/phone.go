package service

import (
	"encoding/json"
	"fmt"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/utils"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type PhoneService struct {
	LoggerSugar                   *zap.SugaredLogger
	PhoneDomainDataBaseRepository output.IPhoneDomainDataBaseRepository
	PhoneDomainCacheRepository    output.IPhoneDomainCacheRepository
}

var PhoneCacheTTL = 10 * time.Minute

const (
	PhoneCacheKeyTypeID = "ID"
)

const (
	PhoneErrorToSaveInCache    = "error to save phone in cache"
	PhoneErrorToGetByIDInCache = "error to get phone by id in cache"
)

func (service *PhoneService) getCacheKey(cacheKeyType, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service *PhoneService) Create(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {
	err := utils.ValidateCodeAreaNumber(phone.CodeArea)
	if err != nil {
		return domain.PhoneDomain{}, nil
	}

	err = utils.ValidatePhoneNumber(phone.PhoneType, phone.Number)
	if err != nil {
		return domain.PhoneDomain{}, err
	}

	save, err := service.PhoneDomainDataBaseRepository.Save(contextControl, phone)
	if err != nil {
		return domain.PhoneDomain{}, err
	}

	hash, _ := json.Marshal(save)
	if err = service.PhoneDomainCacheRepository.Set(contextControl, service.getCacheKey(PhoneCacheKeyTypeID, strconv.FormatInt(phone.ID, 10)),
		string(hash), PhoneCacheTTL); err != nil {
		service.LoggerSugar.Infow(PhoneErrorToSaveInCache, "phone_id", phone.ID)
	}
	return save, nil
}

func (service *PhoneService) GetByID(contextControl domain.ContextControl, ID int64) (domain.PhoneDomain, bool, error) {
	phone, exists, err := service.PhoneDomainDataBaseRepository.GetByID(contextControl, ID)
	if err != nil {
		return domain.PhoneDomain{}, exists, err
	}

	if !exists {
		return domain.PhoneDomain{}, exists, nil
	}
	hash, _ := json.Marshal(phone)
	if err = service.PhoneDomainCacheRepository.Set(contextControl,
		service.getCacheKey(AddressCacheKeyTypeID, strconv.FormatInt(phone.ID, 10)),
		string(hash), PhoneCacheTTL); err != nil {
		service.LoggerSugar.Infow(PhoneErrorToGetByIDInCache, "phone_id", phone.ID)
	}
	return phone, exists, nil
}
