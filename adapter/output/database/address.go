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
	AddressSaveDBError = "failed to save address to postgres"
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

func (c AddressDB) CopyToAddressDomain() (domain.AddressDomain, error) {

	var addressDomain domain.AddressDomain
	if err := copier.Copy(&addressDomain, &c); err != nil {
		return domain.AddressDomain{}, err
	}

	return addressDomain, nil
}

func (cp AddressPostgresDB) Save(contextControl domain.ContextControl, addressDomain domain.AddressDomain) (domain.AddressDomain, error) {

	var addressDB AddressDB
	if err := copier.Copy(&addressDB, &addressDomain); err != nil {
		cp.LoggerSugar.Errorw("error copying address domain to DB struct", "error", err.Error())
		return domain.AddressDomain{}, err
	}

	if err := cp.DB.WithContext(contextControl.Context).
		Create(&addressDB).Error; err != nil {
		cp.LoggerSugar.Errorw(AddressSaveDBError,
			"error", err.Error())
		return domain.AddressDomain{}, err
	}

	addressResult, err := addressDB.CopyToAddressDomain()
	if err != nil {
		cp.LoggerSugar.Errorw("error copying DB struct to address domain", "error", err.Error())
		return domain.AddressDomain{}, err
	}

	return addressResult, nil
}

func (cp AddressPostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
	var addressDB AddressDB

	result := cp.DB.WithContext(contextControl.Context).First(&addressDB, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			cp.LoggerSugar.Infow(AddressNotFound, "address_id", ID)
			return domain.AddressDomain{}, false, nil
		}
		cp.LoggerSugar.Errorw("error getting address by ID from DB", "address_id", ID, "error", result.Error.Error())
		return domain.AddressDomain{}, false, result.Error
	}

	addressResult, err := addressDB.CopyToAddressDomain()
	if err != nil {
		cp.LoggerSugar.Errorw("error copying DB struct to address domain", "error", err.Error())
		return domain.AddressDomain{}, false, err
	}

	return addressResult, true, nil
}
