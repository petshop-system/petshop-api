package database

import (
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PersonPostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

const (
	PersonSaveDBError = "error to save the person into postgres"
)

func NewPersonPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) PersonPostgresDB {
	return PersonPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

func (cp PersonPostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.PersonDomain, error) {
	//TODO implement me
	panic("implement me")
}

type PersonDB struct {
	ID          int64  `gorm:"primaryKey, column:id"`
	Cpf_cnpj    string `gorm:"column:cpf_cnpj"`
	Tipo_pessoa string `gorm:"column:tipo_pessoa"`
}

func (PersonDB) TableName() string {
	return "petshop_api.person"
}

func (c PersonDB) CopyToPersonDomain() domain.PersonDomain {
	return domain.PersonDomain{
		ID:          c.ID,
		Cpf_cnpj:    c.Cpf_cnpj,
		Tipo_pessoa: c.Tipo_pessoa,
	}
}

func (cp PersonPostgresDB) Save(contextControl domain.ContextControl, personDomain domain.PersonDomain) (domain.PersonDomain, error) {

	var personDB PersonDB
	copier.Copy(&personDB, &personDomain)
	personDB.Cpf_cnpj = utils.RemoveNonNumericCharacters(personDB.Cpf_cnpj)

	if err := cp.DB.WithContext(contextControl.Context).
		Create(&personDB).Error; err != nil {
		cp.LoggerSugar.Errorw(PersonSaveDBError,
			"error", err.Error())
		return domain.PersonDomain{}, err
	}

	return personDB.CopyToPersonDomain(), nil
}
