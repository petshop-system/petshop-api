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
	AddressSaveDBError = "error to save the address into postgres "
)

func NewAddressPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) AddressPostgresDB {
	return AddressPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

func (cp AddressPostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.AddressDomain, error) {
	//TODO implement me
	panic("implement me")
}

type AddressDB struct {
	ID         int64  `gorm:"primaryKey, column:id"`
	Logradouro string `gorm:"column:logradouro"`
	Numero     string `gorm:"column:numero"`
}

func (AddressDB) TableName() string {
	return "petshop_api.address"
}

func (c AddressDB) CopyToAddressDomain() domain.AddressDomain {
	return domain.AddressDomain{
		ID:         c.ID,
		Logradouro: c.Logradouro,
		Numero:     c.Numero,
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
