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
	SuccessToCreateCustomer = "user created with success"
	ErrorToCreateCustomer   = "error to create and process the request"
)

type Customer struct {
	CustomerService input.ICustomerService
	LoggerSugar     *zap.SugaredLogger
}

type CustomerRequest struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Phone       map[string]string `json:"phone"` // key: tipo telefone, val: telefone
	Address     string            `json:"address"`
	DateCreated time.Time         `json:"date_created"`
}

type CustomerResponse struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Phone       map[string]string `json:"phone"` // key: tipo telefone, val: telefone
	Address     string            `json:"address"`
	DateCreated time.Time         `json:"date_created"`
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
