package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/input"
	"go.uber.org/zap"
)

const (
	SuccessToCreatePhone = "phone created with success"
	SuccessToGetPhone    = "phone found with success"
	ErrorToCreatePhone   = "error to create and process the request"
	ErrorToGetPhone      = "error retrieving phone by ID" //TODO: Adjust error messages to this model
	PhoneNotFound        = "phone not found"
	PhoneNotFoundMessage = "the phone with id %d wasn't found"
)

type Phone struct {
	PhoneService input.IPhoneService
	LoggerSugar  *zap.SugaredLogger
}

type PhoneRequest struct {
	ID             int64  `json:"id"`
	Number         string `json:"number"`
	CodeAreaNumber string `json:"code_area"`
	PhoneType      string `json:"phone_type"`
}

type PhoneResponse struct {
	ID             int64  `json:"id"`
	Number         string `json:"number"`
	CodeAreaNumber string `json:"code_area"`
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

func (c *Phone) GetByID(w http.ResponseWriter, r *http.Request) {
	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var IDRequest, err = strconv.ParseInt(chi.URLParam(r, "id"), 10, 64) //TODO: I will create a function to streamline this step in an upcoming PR.
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetPhone, "error", err.Error())
		response := objectResponse(ErrorToGetPhone, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	phoneDomain, exists, err := c.PhoneService.GetByID(contextControl, IDRequest)
	if err != nil {
		c.LoggerSugar.Errorw(ErrorToGetPhone, "error", err.Error())
		response := objectResponse(ErrorToGetPhone, err.Error())
		responseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if !exists {
		c.LoggerSugar.Errorw(PhoneNotFound)
		response := objectResponse(PhoneNotFound, fmt.Sprintf(PhoneNotFoundMessage, IDRequest))
		responseReturn(w, http.StatusNotFound, response.Bytes())
		return
	}

	var phoneResponse PhoneResponse
	copier.Copy(&phoneResponse, &phoneDomain)
	response := objectResponse(phoneResponse, SuccessToGetPhone)
	responseReturn(w, http.StatusOK, response.Bytes())
}
