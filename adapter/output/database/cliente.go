package database

import (
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ClientePostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

func NewClientePostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) ClientePostgresDB {
	return ClientePostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

func (cp ClientePostgresDB) GetByID(contextControl domain.ContextControl, ID int64) (domain.ClienteDomain, error) {
	//TODO implement me
	panic("implement me")
}

type ClienteDB struct {
	ID           int64             `gorm:"primaryKey, column:id"`
	Nome         string            `gorm:"column:nome"`
	Telefone     map[string]string `gorm:"type:json, column:telefone"` // key: tipo telefone, val: telefone
	Endereco     string            `gorm:"column:endereco"`
	DataCadastro time.Time         `gorm:"column:data_cadastro"`
}

func (ClienteDB) TableName() string {
	return "petshop_api.cliente"
}

func (c ClienteDB) CopyToClienteDomain() domain.ClienteDomain {
	return domain.ClienteDomain{
		ID:           c.ID,
		Nome:         c.Nome,
		Telefone:     c.Telefone,
		Endereco:     c.Endereco,
		DataCadastro: c.DataCadastro,
	}
}

func (cp ClientePostgresDB) Save(contextControl domain.ContextControl, clienteDomain domain.ClienteDomain) (domain.ClienteDomain, error) {

	var clienteDB ClienteDB
	copier.Copy(&clienteDB, &clienteDomain)

	if err := cp.DB.WithContext(contextControl.Context).
		Create(&clienteDB).Error; err != nil {
		cp.LoggerSugar.Errorw("error to save into postgres",
			"error", err.Error())
		return domain.ClienteDomain{}, err
	}

	return clienteDB.CopyToClienteDomain(), nil
}
