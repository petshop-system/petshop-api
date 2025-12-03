package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/application/port/output"
	"github.com/petshop-system/petshop-api/application/service"
	"github.com/petshop-system/petshop-api/application/utils"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

var pathAddressCreate = "/address/create"

func TestAddress_Create(t *testing.T) {
	t.Run("create address successfully", func(t *testing.T) {
		mockRepo := output.AddressDomainDataBaseRepositoryMock{
			SaveMock: func(ctx domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
				mocked := utils.GetMockAddress()
				mocked.ID = 1
				return mocked, nil
			},
		}
		mockCache := output.AddressDomainCacheRepositoryMock{
			SetMock: func(ctx domain.ContextControl, key string, hash string, expiration time.Duration) error {
				return nil
			},
		}

		addressService := service.AddressService{
			LoggerSugar:                     zap.NewNop().Sugar(),
			AddressDomainDataBaseRepository: mockRepo,
			AddressDomainCacheRepository:    mockCache,
		}
		handler := Address{AddressService: addressService, LoggerSugar: zap.NewNop().Sugar()}

		requestBody := AddressRequest{
			Street:       utils.MockAddressStreet,
			Number:       utils.MockAddressNumber,
			Complement:   utils.MockAddressComplement,
			Neighborhood: utils.MockAddressNeighborhood,
			ZipCode:      utils.MockAddressZipCode,
			City:         utils.MockAddressCity,
			State:        utils.MockAddressState,
			Country:      utils.MockAddressCountry,
		}
		body := new(bytes.Buffer)
		_ = json.NewEncoder(body).Encode(requestBody)

		req := httptest.NewRequest(http.MethodPost, pathAddressCreate, body)
		w := httptest.NewRecorder()
		handler.Create(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run("invalid json in request body", func(t *testing.T) {
		addressService := service.AddressService{}
		handler := Address{AddressService: addressService, LoggerSugar: zap.NewNop().Sugar()}

		body := bytes.NewBufferString(`{"street":123}`)

		req := httptest.NewRequest(http.MethodPost, pathAddressCreate, body)
		w := httptest.NewRecorder()
		handler.Create(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("validation fails - street required", func(t *testing.T) {
		addressService := service.AddressService{}
		handler := Address{AddressService: addressService, LoggerSugar: zap.NewNop().Sugar()}

		requestBody := AddressRequest{
			Street:       "",
			Number:       "12",
			Complement:   "Apto 1",
			Neighborhood: "Copacabana",
			ZipCode:      "12345-678",
			City:         "Rio de Janeiro",
			State:        "RJ",
			Country:      "Brasil",
		}

		body := new(bytes.Buffer)
		_ = json.NewEncoder(body).Encode(requestBody)

		req := httptest.NewRequest(http.MethodPost, pathAddressCreate, body)
		w := httptest.NewRecorder()
		handler.Create(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("database error on save", func(t *testing.T) {
		mockRepo := output.AddressDomainDataBaseRepositoryMock{
			SaveMock: func(ctx domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
				return domain.AddressDomain{}, errors.New("db failure")
			},
		}
		mockCache := output.AddressDomainCacheRepositoryMock{}

		addressService := service.AddressService{
			LoggerSugar:                     zap.NewNop().Sugar(),
			AddressDomainDataBaseRepository: mockRepo,
			AddressDomainCacheRepository:    mockCache,
		}
		handler := Address{AddressService: addressService, LoggerSugar: zap.NewNop().Sugar()}

		requestBody := AddressRequest{
			Street:       "Rua A",
			Number:       "10",
			Complement:   "",
			Neighborhood: "Centro",
			ZipCode:      "12345-678",
			City:         "Rio de Janeiro",
			State:        "RJ",
			Country:      "Brasil",
		}

		body := new(bytes.Buffer)
		_ = json.NewEncoder(body).Encode(requestBody)

		req := httptest.NewRequest(http.MethodPost, pathAddressCreate, body)
		w := httptest.NewRecorder()
		handler.Create(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("cache error after save", func(t *testing.T) {
		mockRepo := output.AddressDomainDataBaseRepositoryMock{
			SaveMock: func(ctx domain.ContextControl, address domain.AddressDomain) (domain.AddressDomain, error) {
				mocked := utils.GetMockAddress()
				mocked.ID = 2
				return mocked, nil
			},
		}
		mockCache := output.AddressDomainCacheRepositoryMock{
			SetMock: func(ctx domain.ContextControl, key string, hash string, expiration time.Duration) error {
				return errors.New("cache failure")
			},
		}

		addressService := service.AddressService{
			LoggerSugar:                     zap.NewNop().Sugar(),
			AddressDomainDataBaseRepository: mockRepo,
			AddressDomainCacheRepository:    mockCache,
		}
		handler := Address{AddressService: addressService, LoggerSugar: zap.NewNop().Sugar()}

		requestBody := AddressRequest{
			Street:       "Rua B",
			Number:       "25",
			Complement:   "",
			Neighborhood: "Zona Sul",
			ZipCode:      "22222-222",
			City:         "Rio",
			State:        "RJ",
			Country:      "Brasil",
		}

		body := new(bytes.Buffer)
		_ = json.NewEncoder(body).Encode(requestBody)

		req := httptest.NewRequest(http.MethodPost, pathAddressCreate, body)
		w := httptest.NewRecorder()
		handler.Create(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
