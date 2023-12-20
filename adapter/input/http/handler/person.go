package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/input"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	SuccessToCreatePerson = "person created with success"
	SuccessToGetPerson    = "person found with success"
	ErrorToCreatePerson   = "error to create and process the request"
	ErrorToGetPerson      = "error to get person by id"
	PersonNotFound        = "person not found"
)

type Person struct {
	PersonService input.IPersonService
	LoggerSugar   *zap.SugaredLogger
}

type PersonRequest struct {
	ID         int64  `json:"id"`
	Document   string `json:"document"`
	PersonType string `json:"person_type"`
}

type PersonResponse struct {
	ID         int64  `json:"id"`
	Document   string `json:"document"`
	PersonType string `json:"person_type"`
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

func (c *Person) GetByID(w http.ResponseWriter, r *http.Request) {
	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var IDRequest, err = strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetPerson, "error", err.Error())
		response := objectResponse(ErrorToGetPerson, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	personDomain, exists, err := c.PersonService.GetByID(contextControl, IDRequest)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetPerson, "error", err.Error())
		response := objectResponse(ErrorToGetPerson, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if !exists {
		c.LoggerSugar.Errorw(PersonNotFound)
		response := objectResponse(PersonNotFound, fmt.Sprintf(PersonNotFound, IDRequest))
		responseReturn(w, http.StatusNotFound, response.Bytes())
		return
	}

	var personResponse PersonResponse
	copier.Copy(&personResponse, &personDomain)
	response := objectResponse(personResponse, SuccessToGetPerson)
	responseReturn(w, http.StatusOK, response.Bytes())
}
