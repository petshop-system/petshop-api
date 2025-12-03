package database

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AddressPostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

const (
	AddressSaveDBError = "error to save the address into postgres"
	AddressNotFound    = "address not found"
)

func NewAddressPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) AddressPostgresDB {
	return AddressPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

type AddressDB struct {
	ID           int64  `gorm:"primaryKey, column:id"`
	Street       string `gorm:"column:street"`
	Number       string `gorm:"column:number"`
	Complement   string `gorm:"column:complement"`
	Neighborhood string `gorm:"column:neighborhood"`
	ZipCode      string `gorm:"column:zip_code"`
	City         string `gorm:"column:city"`
	State        string `gorm:"column:state"`
	Country      string `gorm:"column:country"`
}

func (AddressDB) TableName() string {
	return "petshop_api.address"
}

func (c AddressDB) CopyToAddressDomain() domain.AddressDomain {

	var addressDomain domain.AddressDomain
	copier.Copy(&addressDomain, &c)

	return addressDomain
}

func (cp AddressPostgresDB) Save(contextControl domain.ContextControl, addressDomain domain.AddressDomain) (domain.AddressDomain, error) {

	var addressDB AddressDB
	copier.Copy(&addressDB, &addressDomain)

	if err := cp.DB.WithContext(contextControl.Context).
		Create(&addressDB).Error; err != nil {
		cp.LoggerSugar.Errorw(AddressSaveDBError,
			"error", err.Error())
		return domain.AddressDomain{}, err
	}

	return addressDB.CopyToAddressDomain(), nil
}

func (cp AddressPostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
	var addressDB AddressDB

	result := cp.DB.WithContext(contextControl.Context).First(&addressDB, ID)
	if result.RowsAffected == 0 {
		cp.LoggerSugar.Errorw(AddressNotFound)
		return domain.AddressDomain{}, false, errors.New(AddressNotFound)
	}

	return addressDB.CopyToAddressDomain(), true, nil
}
