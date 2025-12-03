package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/input"
	"go.uber.org/zap"
)

const (
	SuccessToCreateCustomer       = "user created with success"
	ErrorToCreateCustomer         = "error to create and process the request"
	ErrorValidateCreateCustomer   = "validation got some mistakes"
	SuccessValidateCreateCustomer = "success to validate create customer"
)

type Customer struct {
	CustomerService input.ICustomerService
	LoggerSugar     *zap.SugaredLogger
}

type CustomerRequest struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Document   string `json:"document"`
	PersonType string `json:"person_type"`
	ContractID int64  `json:"contract_id"`
	AddressID  int64  `json:"address_id"`
}

type CustomerResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Document   string `json:"document"`
	PersonType string `json:"person_type"`
	ContractID int64  `json:"contract_id"`
	AddressID  int64  `json:"address_id"`
}

func (c *Customer) Create(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var customerRequest CustomerRequest
	json.NewDecoder(r.Body).Decode(&customerRequest)

	var customerDomain domain.CustomerDomain
	copier.Copy(&customerDomain, &customerRequest)

	customerDomain, err := c.CustomerService.Create(contextControl, customerDomain)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToCreateCustomer, "error", err.Error())
		response := objectResponse(ErrorToCreateCustomer, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var customerResponse CustomerResponse
	copier.Copy(&customerResponse, &customerDomain)
	response := objectResponse(customerResponse, SuccessToCreateCustomer)
	responseReturn(w, http.StatusCreated, response.Bytes())
}

func (c *Customer) ValidateCreate(w http.ResponseWriter, r *http.Request) {

	var customerRequest CustomerRequest
	json.NewDecoder(r.Body).Decode(&customerRequest)

	var customerDomain domain.CustomerDomain
	copier.Copy(&customerDomain, &customerRequest)

	if err := c.CustomerService.ValidateCreate(customerDomain); err != nil {
		c.LoggerSugar.Errorw(ErrorValidateCreateCustomer, "error", err.Error())
		response := objectResponse(ErrorValidateCreateCustomer, err.Error())
		responseReturn(w, http.StatusBadRequest, response.Bytes())
		return
	}

	var customerResponse CustomerResponse
	copier.Copy(&customerResponse, &customerDomain)
	response := objectResponse(customerResponse, SuccessValidateCreateCustomer)
	responseReturn(w, http.StatusOK, response.Bytes())
}
