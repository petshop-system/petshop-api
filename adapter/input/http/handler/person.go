package handler

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/input"
	"go.uber.org/zap"
	"net/http"
)

const (
	SuccessToCreatePerson = "person created with success"
	ErrorToCreatePerson   = "error to create and process the request"
)

type Person struct {
	PersonService input.IPersonService
	LoggerSugar   *zap.SugaredLogger
}

type PersonRequest struct {
	ID          int64  `json:"id"`
	Cpf_cnpj    string `json:"cpf_cnpj"`
	Tipo_pessoa string `json:"tipo_pessoa"`
}

type PersonResponse struct {
	ID          int64  `json:"id"`
	Cpf_cnpj    string `json:"cpf_cnpj"`
	Tipo_pessoa string `json:"tipo_pessoa"`
}

func (c *Person) Create(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var personRequest PersonRequest
	json.NewDecoder(r.Body).Decode(&personRequest)

	var personDomain domain.PersonDomain
	copier.Copy(&personDomain, &personRequest)

	personDomain, err := c.PersonService.Create(contextControl, personDomain)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToCreatePerson, "error", err.Error())
		response := objectResponse(ErrorToCreatePerson, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var personResponse PersonResponse
	copier.Copy(&personResponse, &personDomain)
	response := objectResponse(personResponse, SuccessToCreatePerson)
	responseReturn(w, http.StatusCreated, response.Bytes())
}
