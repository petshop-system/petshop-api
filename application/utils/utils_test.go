package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCpf(t *testing.T) {
	tests := []struct {
		Name          string
		InputCpf      string
		ExpectedError error
	}{
		{
			Name:          "valid CPF format",
			InputCpf:      "013.405.400-88",
			ExpectedError: nil,
		},
		{
			Name:          "invalid CPF format",
			InputCpf:      "0.45.40-8",
			ExpectedError: errors.New(ErrorInvalidCPFLength),
		},
		{
			Name:          "invalid CPF with all digits equal",
			InputCpf:      "111.111.111-11",
			ExpectedError: errors.New(ErrorAllDigitsEqualCPF),
		},
		{
			Name:          "valid CPF with extra characters",
			InputCpf:      "123.456.789-09-X",
			ExpectedError: errors.New(ErrorInvalidCPFLength),
		},
		{
			Name:          "invalid CPF with incorrect first verification digit",
			InputCpf:      "076.164.346-16",
			ExpectedError: errors.New(ErrorFirstVerificationCPF),
		},
		{
			Name:          "invalid CPF with incorrect second verification digits",
			InputCpf:      "529.982.247-27",
			ExpectedError: errors.New(ErrorSecondVerificationCPF),
		},
		{
			Name:          "invalid CPF with incorrect character conversion",
			InputCpf:      "12A.345.678-90",
			ExpectedError: errors.New(ErrorInvalidCPFLength),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := ValidateCpf(test.InputCpf)
			assert.Equal(t, err, test.ExpectedError)
		})
	}
}

func TestValidateCnpj(t *testing.T) {
	tests := []struct {
		Name          string
		InputCnpj     string
		ExpectedError error
	}{
		{
			Name:          "valid CNPJ format",
			InputCnpj:     "79.626.068/0001-30",
			ExpectedError: nil,
		},
		{
			Name:          "invalid CNPJ format",
			InputCnpj:     "79626",
			ExpectedError: errors.New(ErrorInvalidCNPJLength),
		},
		{
			Name:          "invalid CNPJ with all digits equal",
			InputCnpj:     "11.111.111/1111-11",
			ExpectedError: errors.New(ErrorAllDigitsEqualCNPJ),
		},
		{
			Name:          "invalid CNPJ with extra characters",
			InputCnpj:     "123.456.789-09-X",
			ExpectedError: errors.New(ErrorInvalidCNPJLength),
		},
		{
			Name:          "invalid CNPJ with incorrect first verification digit",
			InputCnpj:     "79.626.068/0001-00",
			ExpectedError: errors.New(ErrorFirstVerificationCNPJ),
		},
		{
			Name:          "invalid CNPJ with incorrect second verification digits",
			InputCnpj:     "79.626.068/0001-39",
			ExpectedError: errors.New(ErrorSecondVerificationCNPJ),
		},
		{
			Name:          "invalid CNPJ with incorrect character conversion",
			InputCnpj:     "7K.626.068/0001-30",
			ExpectedError: errors.New(ErrorInvalidCNPJLength),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := ValidateCnpj(test.InputCnpj)
			assert.Equal(t, err, test.ExpectedError)
		})
	}
}

func TestRemoveNonAlphaNumericCharacters(t *testing.T) {
	testCases := []struct {
		Name           string
		InputDocument  string
		ExpectedResult string
	}{
		{
			Name:           "given a cnpj with mask must remove non Characters",
			InputDocument:  "123.456-789/0001-23",
			ExpectedResult: "123456789000123",
		},
		{
			Name:           "given a string with latter and number must remove non alpha Characters",
			InputDocument:  "abc.123",
			ExpectedResult: "abc123",
		},
		{
			Name:           "given a string of numbers must return only the given numbers",
			InputDocument:  "4567",
			ExpectedResult: "4567",
		},
		{
			Name:           "given a string of special chars must return an empty string",
			InputDocument:  "!@#$%^&*()_-+=/\\|?",
			ExpectedResult: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			result := RemoveNonAlphaNumericCharacters(test.InputDocument)
			assert.Equal(t, test.ExpectedResult, result)
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	testCases := []struct {
		Name          string
		PhoneNumber   string
		PhoneType     string
		ExpectedError error
	}{
		{
			Name:          "valid mobile phone format",
			PhoneNumber:   "99919-2111",
			PhoneType:     "mobile_phone",
			ExpectedError: nil,
		},
		{
			Name:          "invalid mobile phone format",
			PhoneNumber:   "1234",
			PhoneType:     "mobile_phone",
			ExpectedError: errors.New(ErrorInvalidMobilePhoneLength),
		},
		{
			Name:          "valid landline phone format",
			PhoneNumber:   "3212-2111",
			PhoneType:     "landline_phone",
			ExpectedError: nil,
		},
		{
			Name:          "invalid landline phone format",
			PhoneNumber:   "1234",
			PhoneType:     "landline_phone",
			ExpectedError: errors.New(ErrorInvalidLandLinePhoneLength),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			err := ValidatePhoneNumber(test.PhoneType, test.PhoneNumber)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}

func TestValidateCodeAreaNumber(t *testing.T) {
	testCases := []struct {
		Name          string
		AreaCode      string
		ExpectedError error
	}{
		{
			Name:          "valid area code format",
			AreaCode:      "32",
			ExpectedError: nil,
		},
		{
			Name:          "invalid area code format",
			AreaCode:      "29",
			ExpectedError: errors.New(ErrorAreaCodeVerification),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			err := ValidateCodeAreaNumber(test.AreaCode)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}
