package database

import (
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PhonePostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

const (
	PhoneSaveError      = "error to save the phone into postgres"
	PhoneGetByIDDBError = "error to get a phone by id"
	PhoneNotFound       = "phone not found"
)

func NewPhonePostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) PhonePostgresDB {
	return PhonePostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

type PhoneDB struct {
	ID        int64  `gorm:"primaryKey, column:id"`
	Number    string `gorm:"column:number"`
	CodeArea  string `gorm:"column:code_area"`
	PhoneType string `gorm:"column:phone_type"`
}

func (PhoneDB) TableName() string {
	return "petshop_api.phone"
}

func (c PhoneDB) CopyToPhoneDomain() domain.PhoneDomain {
	return domain.PhoneDomain{
		ID:        c.ID,
		Number:    c.Number,
		CodeArea:  c.CodeArea,
		PhoneType: c.PhoneType,
	}
}

func (cp PhonePostgresDB) Save(contextControl domain.ContextControl, phoneDomain domain.PhoneDomain) (domain.PhoneDomain, error) {

	var phoneDB PhoneDB
	copier.Copy(&phoneDB, &phoneDomain)
	phoneDB.Number = utils.RemoveNonAlphaNumericCharacters(phoneDB.Number)

	if err := cp.DB.WithContext(contextControl.Context).Create(&phoneDB).Error; err != nil {
		cp.LoggerSugar.Errorw(PhoneSaveError, "error", err.Error())
		return domain.PhoneDomain{}, err
	}
	return phoneDB.CopyToPhoneDomain(), nil
}

func (cp PhonePostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.PhoneDomain, bool, error) {

	var phoneDB PhoneDB

	result := cp.DB.WithContext(contextControl.Context).First(&phoneDB, ID)
	if result.RowsAffected == 0 {
		cp.LoggerSugar.Errorw(PhoneNotFound)
		return domain.PhoneDomain{}, false, nil
	}
	return phoneDB.CopyToPhoneDomain(), true, nil
}
