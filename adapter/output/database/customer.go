package database

import (
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type CustomerPostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

func NewCustomerPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) CustomerPostgresDB {
	return CustomerPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

func (cp CustomerPostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.CustomerDomain, error) {
	//TODO implement me
	panic("implement me")
}

type CustomerDB struct {
	ID          int64             `gorm:"primaryKey, column:id"`
	Name        string            `gorm:"column:name"`
	Phone       map[string]string `gorm:"type:json, column:phone"` // key: tipo telefone, val: telefone
	Address     string            `gorm:"column:address"`
	DataCreated time.Time         `gorm:"column:date_created"`
}

func (CustomerDB) TableName() string {
	return "petshop_api.customer"
}

func (c CustomerDB) CopyToCustomerDomain() domain.CustomerDomain {
	return domain.CustomerDomain{
		ID:          c.ID,
		Name:        c.Name,
		DateCreated: c.DataCreated,
	}
}

func (cp CustomerPostgresDB) Save(contextControl domain.ContextControl, customerDomain domain.CustomerDomain) (domain.CustomerDomain, error) {

	var customerDB CustomerDB
	copier.Copy(&customerDB, &customerDomain)

	if err := cp.DB.WithContext(contextControl.Context).
		Create(&customerDB).Error; err != nil {
		cp.LoggerSugar.Errorw("error to save into postgres",
			"error", err.Error())
		return domain.CustomerDomain{}, err
	}

	return customerDB.CopyToCustomerDomain(), nil
}
