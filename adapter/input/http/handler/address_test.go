package handler

import (
	"bytes"
	"encoding/json"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAddressService struct {
	CreateFunc  func(ctx domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error)
	GetByIDFunc func(ctx domain.ContextControl, ID int64) (domain.AddressDomain, bool, error)
}

func (m *MockAddressService) Create(ctx domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
	return m.CreateFunc(ctx, address)
}

func (m *MockAddressService) GetByID(ctx domain.ContextControl, ID int64) (domain.AddressDomain, bool, error) {
	return m.GetByIDFunc(ctx, ID)
}

func TestAddress_Create(t *testing.T) {
	tests := []struct {
		name           string
		addressRequest AddressRequest
		mockService    *MockAddressService
		expectedCode   int
		expectedBody   string
	}{
		{
			name: "Test Successful - Create",
			addressRequest: AddressRequest{
				ID:           1,
				Street:       "Rua do Mairon",
				Number:       "5",
				Complement:   "Casa 7",
				Block:        "",
				Neighborhood: "Flamengo",
				ZipCode:      "12345-678",
				City:         "Rio de Janeiro 1",
				State:        "RJ",
				Country:      "Brasil",
			},
			mockService: &MockAddressService{
				CreateFunc: func(ctx domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					return address, nil
				},
			},
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"address created with success","result":{"id":1,"street":"Rua do Mairon","number":"5","complement":"Casa 7","block":"","neighborhood":"Flamengo","zip_code":"12345-678","city":"Rio de Janeiro 1","state":"RJ","country":"Brasil"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := zap.NewDevelopment()
			sugar := logger.Sugar()

			handler := &Address{
				AddressService: tt.mockService,
				LoggerSugar:    sugar,
			}

			body, _ := json.Marshal(tt.addressRequest)
			request := httptest.NewRequest(http.MethodPost, "/address/search", bytes.NewBuffer(body))
			response := httptest.NewRecorder()
			handler.Create(response, request)

			assert.Equal(t, tt.expectedCode, response.Code, "Expected status code did not match")

			assert.Equal(t, tt.expectedBody, response.Body.String()[:len(tt.expectedBody)], "Expected body did not match")
		})
	}
}
