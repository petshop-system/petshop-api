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
	PersonSaveDBError    = "error to save the person into postgres"
	PersonGetByIdDBError = "error to get a person by id"
	PersonNotFound       = "person not found"
)

func NewPersonPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) PersonPostgresDB {
	return PersonPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

type PersonDB struct {
	ID          int64  `gorm:"primaryKey, column:id"`
	Document    string `gorm:"column:document"`
	Person_type string `gorm:"column:person_type"`
}

func (PersonDB) TableName() string {
	return "petshop_api.person"
}

func (c PersonDB) CopyToPersonDomain() domain.PersonDomain {
	return domain.PersonDomain{
		ID:          c.ID,
		Document:    c.Document,
		Person_type: c.Person_type,
	}
}

func (cp PersonPostgresDB) Save(contextControl domain.ContextControl, personDomain domain.PersonDomain) (domain.PersonDomain, error) {

	var personDB PersonDB
	copier.Copy(&personDB, &personDomain)
	personDB.Document = utils.RemoveNonAlphaNumericCharacters(personDB.Document)

	if err := cp.DB.WithContext(contextControl.Context).Create(&personDB).Error; err != nil {
		cp.LoggerSugar.Errorw(PersonSaveDBError, "error", err.Error())
		return domain.PersonDomain{}, err
	}
	return personDB.CopyToPersonDomain(), nil
}

func (cp PersonPostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.PersonDomain, bool, error) {
	var personDB PersonDB

	result := cp.DB.WithContext(contextControl.Context).First(&personDB, ID)
	if result.RowsAffected == 0 {
		cp.LoggerSugar.Errorw(PersonNotFound)
		return domain.PersonDomain{}, false, nil
	}
	return personDB.CopyToPersonDomain(), true, nil
}
