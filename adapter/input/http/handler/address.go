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
	SuccessToCreateAddress = "address created with success"
	ErrorToCreateAddress   = "error to create and process the request"
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
