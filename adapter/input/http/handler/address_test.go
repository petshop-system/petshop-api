package handler

import (
	"bytes"
	"encoding/json"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAddress_Create(t *testing.T) {

	tests := []struct {
		name            string
		addressRequest  AddressRequest
		mockRepository  output.IAddressDomainDataBaseRepository
		cacheRepository output.IAddressDomainCacheRepository
		expectedCode    int
		expectedBody    string
	}{
		{
			name: "Test Successful - Create",
			addressRequest: AddressRequest{
				Street:       "Rua do Mairon",
				Number:       "5",
				Complement:   "Casa 7",
				Neighborhood: "Flamengo",
				ZipCode:      "12345-678",
				City:         "Rio de Janeiro",
				State:        "RJ",
				Country:      "Brasil",
			},
			mockRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					return domain.AddressDomain{
						Street:       "Rua do Mairon",
						Number:       "5",
						Complement:   "Casa 7",
						Neighborhood: "Flamengo",
						ZipCode:      "12345-678",
						City:         "Rio de Janeiro",
						State:        "RJ",
						Country:      "Brasil",
					}, nil
				},
			},
			cacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"address created with success","result":{"id":0,"street":"Rua do Mairon","number":"5","complement":"Casa 7","neighborhood":"Flamengo","zip_code":"12345-678","city":"Rio de Janeiro","state":"RJ","country":"Brasil"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger, _ := zap.NewDevelopment()
			sugar := logger.Sugar()

			addressService := service.AddressService{
				LoggerSugar:                     sugar,
				AddressDomainDataBaseRepository: tt.mockRepository,
				AddressDomainCacheRepository:    tt.cacheRepository,
			}

			addressHandler := Address{
				AddressService: addressService,
				LoggerSugar:    zap.NewExample().Sugar(),
			}

			AddressRequestJson, err := json.Marshal(tt.addressRequest)
			if err != nil {
				t.Fatal(err)
			}

			request, err := http.NewRequest(http.MethodPost, "/address/create", bytes.NewBuffer(AddressRequestJson))
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()
			handler := http.HandlerFunc(addressHandler.Create)
			handler.ServeHTTP(response, request)

			assert.Equal(t, tt.expectedCode, response.Code, "Unexpected response code")

			var gotBody map[string]interface{}
			err = json.Unmarshal(response.Body.Bytes(), &gotBody)
			assert.NoError(t, err, "Error unmarshalling response body")

			delete(gotBody, "date")

			var expectedBody map[string]interface{}
			err = json.Unmarshal([]byte(tt.expectedBody), &expectedBody)
			assert.NoError(t, err, "Error unmarshalling expected body")

			assert.Equal(t, expectedBody, gotBody, "Unexpected response body")
		})
	}
}
