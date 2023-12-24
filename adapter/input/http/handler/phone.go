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
	SuccessToCreatePhone = "phone created with success"
	ErrorToCreatePhone   = "error to create and process the request"
)

type Phone struct {
	PhoneService input.IPhoneService
	LoggerSugar  *zap.SugaredLogger
}

type PhoneRequest struct {
	ID             int64  `json:"id"`
	number         string `json:"number"`
	CodeAreaNumber string `json:"code_area_number"`
	PhoneType      string `json:"phone_type"`
}

type PhoneResponse struct {
	ID             int64  `json:"id"`
	number         string `json:"number"`
	CodeAreaNumber string `json:"code_area_number"`
	PhoneType      string `json:"phone_type"`
}

func (c *Phone) Create(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var phoneRequest PhoneRequest
	json.NewDecoder(r.Body).Decode(&phoneRequest)

	var phoneDomain domain.PhoneDomain
	copier.Copy(&phoneDomain, &phoneRequest)

	phoneDomain, err := c.PhoneService.Create(contextControl, phoneDomain)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToCreatePhone, "error", err.Error())
		response := objectResponse(ErrorToCreatePhone, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var phoneResponse PhoneResponse
	copier.Copy(&phoneResponse, &phoneDomain)
	response := objectResponse(phoneResponse, SuccessToCreatePhone)
	responseReturn(w, http.StatusCreated, response.Bytes())
}
