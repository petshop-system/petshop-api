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
	SuccessToCreateAddress = "address created with success"
	SuccessToGetAddress    = "address found with success"
	ErrorToCreateAddress   = "error to create and process the request"
	ErrorToGetAddress      = "error to get and address by id"
	AddressNotFound        = "address not found"
	AddressNotFoundMessage = "the address with id %d wasn't found"
)

type Address struct {
	AddressService input.IAddressService
	LoggerSugar    *zap.SugaredLogger
}

type AddressRequest struct {
	ID         int64  `json:"id"`
	Logradouro string `json:"logradouro"`
	Numero     string `json:"numero"`
}

type AddressResponse struct {
	ID         int64  `json:"id"`
	Logradouro string `json:"logradouro"`
	Numero     string `json:"numero"`
}

func (c *Address) Create(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var addressRequest AddressRequest
	json.NewDecoder(r.Body).Decode(&addressRequest)

	var addressDomain domain.AddressDomain
	copier.Copy(&addressDomain, &addressRequest)

	addressDomain, err := c.AddressService.Create(contextControl, addressDomain)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToCreateAddress, "error", err.Error())
		response := objectResponse(ErrorToCreateAddress, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var addressResponse AddressResponse
	copier.Copy(&addressResponse, &addressDomain)
	response := objectResponse(addressResponse, SuccessToCreateAddress)
	responseReturn(w, http.StatusCreated, response.Bytes())
}

func (c *Address) GetByID(w http.ResponseWriter, r *http.Request) {
	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var IDRequest, err = strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetAddress, "error", err.Error())
		response := objectResponse(ErrorToGetAddress, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	addressDomain, exists, err := c.AddressService.GetByID(contextControl, IDRequest)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetAddress, "error", err.Error())
		response := objectResponse(ErrorToGetAddress, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if !exists {
		c.LoggerSugar.Errorw(AddressNotFound)
		response := objectResponse(AddressNotFound, fmt.Sprintf(AddressNotFoundMessage, IDRequest))
		responseReturn(w, http.StatusNotFound, response.Bytes())
		return
	}

	var addressResponse AddressResponse
	copier.Copy(&addressResponse, &addressDomain)
	response := objectResponse(addressResponse, SuccessToGetAddress)
	responseReturn(w, http.StatusOK, response.Bytes())
}
