package handler

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/input"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	SuccessToCreateCliente = "usuário cadastrado com sucesso"
	ErrorToCreateCliente   = "error to create and process the request"
)

type Cliente struct {
	CustomerService input.ICustomerService
	LoggerSugar     *zap.SugaredLogger
}

type ClienteRequest struct {
	ID           int64             `json:"id"`
	Nome         string            `json:"nome"`
	Telefone     map[string]string `json:"telefone"` // key: tipo telefone, val: telefone
	Endereco     string            `json:"endereco"`
	DataCadastro time.Time         `json:"data_cadastro"`
}

type ClienteResponse struct {
	ID           int64             `json:"id"`
	Nome         string            `json:"nome"`
	Telefone     map[string]string `json:"telefone"` // key: tipo telefone, val: telefone
	Endereco     string            `json:"endereco"`
	DataCadastro time.Time         `json:"data_cadastro"`
}

func (c *Cliente) Create(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var clienteRequest ClienteRequest
	json.NewDecoder(r.Body).Decode(&clienteRequest)

	var clienteDomain domain.ClienteDomain
	copier.Copy(&clienteDomain, &clienteRequest)

	clienteDomain, err := c.CustomerService.Create(contextControl, clienteDomain)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToCreateCliente, "error", err.Error())
		response := objectResponse(ErrorToCreateCliente, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var clienteResponse ClienteResponse
	copier.Copy(&clienteResponse, &clienteDomain)
	response := objectResponse(clienteResponse, SuccessToCreateCliente)
	responseReturn(w, http.StatusCreated, response.Bytes())
}
