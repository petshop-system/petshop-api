package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/service"
	"github.com/petshop-system/petshop-api/application/utils"
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
				Street:       utils.MockAddressStreet,
				Number:       utils.MockAddressNumber,
				Complement:   utils.MockAddressComplement,
				Neighborhood: utils.MockAddressNeighborhood,
				ZipCode:      utils.MockAddressZipCode,
				City:         utils.MockAddressCity,
				State:        utils.MockAddressState,
				Country:      utils.MockAddressCountry,
			},
			mockRepository: output.AddressDomainDataBaseRepositoryMock{
				SaveMock: func(contextControl domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
					address = utils.GetMockAddress()
					address.ID = 0

					return address, nil
				},
			},
			cacheRepository: output.AddressDomainCacheRepositoryMock{
				SetMock: func(contextControl domain.ContextControl, key string, hash string, expirationTime time.Duration) error {
					return nil
				},
			},
			expectedCode: http.StatusCreated,
			expectedBody: fmt.Sprintf(`{"message":"address created with success","result":{"id":0,"street":"%s","number":"%s","complement":"%s","neighborhood":"%s","zip_code":"%s","city":"%s","state":"%s","country":"%s"}}`,
				utils.MockAddressStreet, utils.MockAddressNumber, utils.MockAddressComplement, utils.MockAddressNeighborhood, utils.MockAddressZipCode, utils.MockAddressCity, utils.MockAddressState, utils.MockAddressCountry),
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
