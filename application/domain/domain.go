package domain

import "time"

type ClienteDomain struct {
	ID           int64
	Nome         string
	Telefone     map[string]string // key: tipo telefone, val: telefone
	Endereco     string
	DataCadastro time.Time
}
