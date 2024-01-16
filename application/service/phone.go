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
	LandLinePhone       = "landline_phone"
	MobilePhone         = "mobile_phone"
)

const (
	PhoneErrorToSaveInCache         = "error to save phone in cache"
	PhoneErrorToGetByIDInCache      = "error to get phone by id in cache"
	ErrorInvalidMobilePhoneLength   = "invalid Mobile Phone length error"
	ErrorInvalidLandLinePhoneLength = "invalid Land Line Phone length error"
	InvalidTypeOfPhone              = "invalid type of phone"
)

func (service *PhoneService) getCacheKey(cacheKeyType, value string) string {
	return fmt.Sprintf("%s.%s", cacheKeyType, value)
}

func (service *PhoneService) Create(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error) {

	err := service.ValidatePhone(phone)
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

func (service *PhoneService) ValidatePhone(phone domain.PhoneDomain) error {
	if _, err := utils.ValidateCodeAreaNumber(phone.CodeArea); err != nil {
		return err
	}

	clearPhone := utils.RemoveNonAlphaNumericCharacters(phone.Number)
	verification := func(phoneLen int, phoneTypeVerification, errorMessageVerification string) error {
		if len(clearPhone) != phoneLen {
			return fmt.Errorf(errorMessageVerification)
		}
		return nil
	}

	switch phone.PhoneType {
	case LandLinePhone:
		landLinePhoneLen := 8
		if err := verification(landLinePhoneLen, phone.PhoneType, ErrorInvalidLandLinePhoneLength); err != nil {
			return err
		}
	case MobilePhone:
		mobilePhoneLen := 9
		if err := verification(mobilePhoneLen, phone.PhoneType, ErrorInvalidMobilePhoneLength); err != nil {
			return err
		}
	default:
		return fmt.Errorf(InvalidTypeOfPhone)
	}
	return nil
}
