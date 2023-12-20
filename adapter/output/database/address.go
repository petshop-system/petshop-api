package database

import (
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
	AddressSaveDBError    = "error to save the address into postgres "
	AddressGetByIdDBError = "error to get an address by id"
	AddressNotFound       = "address not found"
)

func NewAddressPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) AddressPostgresDB {
	return AddressPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

type AddressDB struct {
	ID     int64  `gorm:"primaryKey, column:id"`
	Street string `gorm:"column:street"`
	Number string `gorm:"column:number"`
}

func (AddressDB) TableName() string {
	return "petshop_api.address"
}

func (c AddressDB) CopyToAddressDomain() domain.AddressDomain {
	return domain.AddressDomain{
		ID:     c.ID,
		Street: c.Street,
		Number: c.Number,
	}
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
		return domain.AddressDomain{}, false, nil
	}

	return addressDB.CopyToAddressDomain(), true, nil
}
