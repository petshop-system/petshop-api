package domain

import "time"

type ClienteDomain struct {
	ID           int64
	Nome         string
	DataCadastro time.Time
}

type TelefoneDomain struct {
	ID           int64
	Numero       string
	DDD          string
	TipoTelefone string
}

type EspecieDomain struct {
	ID   int64
	Nome string
}

type RacaDomain struct {
	ID   int64
	Nome string
}

type AddressDomain struct {
	ID         int64
	Logradouro string
	Numero     string
}
