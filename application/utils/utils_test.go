package utils

import (
	"errors"

	"fmt"
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
			Name:          "WithValidCPF_ReturnsNoError",
			InputCpf:      "013.405.400-88",
			ExpectedError: nil,
		},
		{
			Name:          "WithInvalidCPFLength_ReturnsError",
			InputCpf:      "0.45.40-8",
			ExpectedError: errors.New(ErrorInvalidCPFLength),
		},
		{
			Name:          "WithAllDigitsEqual_ReturnsError",
			InputCpf:      "111.111.111-11",
			ExpectedError: errors.New(ErrorAllDigitsEqualCPF),
		},
		{
			Name:          "WithExtraCharacters_ReturnsError",
			InputCpf:      "123.456.789-09-X",
			ExpectedError: errors.New(ErrorInvalidCPFLength),
		},
		{
			Name:          "WithIncorrectFirstVerificationDigit_ReturnsError",
			InputCpf:      "076.164.346-16",
			ExpectedError: errors.New(ErrorFirstVerificationCPF),
		},
		{
			Name:          "WithIncorrectSecondVerificationDigit_ReturnsError",
			InputCpf:      "529.982.247-27",
			ExpectedError: errors.New(ErrorSecondVerificationCPF),
		},
		{
			Name:          "WithInvalidCharacter_ReturnsError",
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
			Name:          "WithValidCNPJ_ReturnsNoError",
			InputCnpj:     "79.626.068/0001-30",
			ExpectedError: nil,
		},
		{
			Name:          "WithInvalidCNPJLength_ReturnsError",
			InputCnpj:     "79626",
			ExpectedError: errors.New(ErrorInvalidCNPJLength),
		},
		{
			Name:          "WithAllDigitsEqual_ReturnsError",
			InputCnpj:     "11.111.111/1111-11",
			ExpectedError: errors.New(ErrorAllDigitsEqualCNPJ),
		},
		{
			Name:          "WithExtraCharacters_ReturnsError",
			InputCnpj:     "123.456.789-09-X",
			ExpectedError: errors.New(ErrorInvalidCNPJLength),
		},
		{
			Name:          "WithIncorrectFirstVerificationDigit_ReturnsError",
			InputCnpj:     "79.626.068/0001-00",
			ExpectedError: errors.New(ErrorFirstVerificationCNPJ),
		},
		{
			Name:          "WithIncorrectSecondVerificationDigit_ReturnsError",
			InputCnpj:     "79.626.068/0001-39",
			ExpectedError: errors.New(ErrorSecondVerificationCNPJ),
		},
		{
			Name:          "WithInvalidCharacter_ReturnsError",
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
			Name:           "WithCNPJMask_RemovesNonAlphanumeric",
			InputDocument:  "123.456-789/0001-23",
			ExpectedResult: "123456789000123",
		},
		{
			Name:           "WithMixedLettersAndNumbers_RemovesSpecialChars",
			InputDocument:  "abc.123",
			ExpectedResult: "abc123",
		},
		{
			Name:           "WithOnlyNumbers_ReturnsUnchanged",
			InputDocument:  "4567",
			ExpectedResult: "4567",
		},
		{
			Name:           "WithOnlySpecialChars_ReturnsEmpty",
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

func TestValidateCodeAreaNumber(t *testing.T) {
	testCases := []struct {
		Name           string
		AreaCode       string
		ExpectedResult string
		ExpectedError  error
	}{
		{
			Name:           "WithValidAreaCode_ReturnsNoError",
			AreaCode:       "32",
			ExpectedResult: "32",
			ExpectedError:  nil,
		},
		{
			Name:           "WithInvalidAreaCode_ReturnsError",
			AreaCode:       "293",
			ExpectedResult: "",
			ExpectedError:  fmt.Errorf(ErrorAreaCodeVerification),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			result, err := ValidateCodeAreaNumber(test.AreaCode)

			assert.Equal(t, test.ExpectedResult, result)

			if test.ExpectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.ExpectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
