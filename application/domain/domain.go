package domain

import (
	"time"
)

type ClienteDomain struct {
	ID           int64
	Nome         string
	DataCadastro time.Time
	Endereco     EnderecoDomain
}

type TelefoneDomain struct {
	ID           int64
	Numero       string
	DDD          string
	TipoTelefone string
	Cliente      ClienteDomain
}

type EnderecoDomain struct {
	ID         int64
	Logradouro string
	Numero     string
}

type Especie struct {
	ID   int64
	Nome string
}

type Raca struct {
	ID      int64
	Nome    string
	Especie Especie
}
