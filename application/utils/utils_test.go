package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCpf(t *testing.T) {
	tests := []struct {
		Name          string
		InputCpf      string
		ExpectedError string
	}{
		{
			Name:          "valid CPF format",
			InputCpf:      "123.456.789-09",
			ExpectedError: "",
		},
		{
			Name:          "valid CPF without formatting",
			InputCpf:      "12345678909",
			ExpectedError: "",
		},
		{
			Name:          "invalid CPF format",
			InputCpf:      "123-456-789",
			ExpectedError: "invalid CPF format",
		},
		{
			Name:          "invalid CPF with all digits equal",
			InputCpf:      "111.111.111-11",
			ExpectedError: "invalid CPF because all digits are equal",
		},
		{
			Name:          "valid CPF with extra characters",
			InputCpf:      "123.456.789-09-X",
			ExpectedError: "invalid CPF format",
		},
		{
			Name:          "invalid CPF with incorrect first verification digit",
			InputCpf:      "529.982.247-26",
			ExpectedError: "error in the second verification of the CPF",
		},
		{
			Name:          "invalid CPF with incorrect second verification digits",
			InputCpf:      "529.982.247-26",
			ExpectedError: "error in the second verification of the CPF",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := ValidateCpf(test.InputCpf)
			if test.ExpectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.ExpectedError)
			}
		})
	}
}

func TestIsValidCpfFormat(t *testing.T) {
	tests := []struct {
		Name           string
		InputCpf       string
		ExpectedResult bool
	}{
		{
			Name:           "valid CPF format with hyphens and dots",
			InputCpf:       "123.456.789-09",
			ExpectedResult: true,
		},
		{
			Name:           "valid CPF format without hyphens and dots",
			InputCpf:       "12345678909",
			ExpectedResult: true,
		},
		{
			Name:           "invalid CPF format with additional characters",
			InputCpf:       "123-456-789-X",
			ExpectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := isValidCpfFormat(test.InputCpf)
			assert.Equal(t, test.ExpectedResult, result)
		})
	}
}

func TestCleanCpf(t *testing.T) {
	tests := []struct {
		Name        string
		InputCpf    string
		ExpectedCpf string
	}{
		{
			Name:        "clean valid CPF format with hyphens and dots",
			InputCpf:    "123.456.789-09",
			ExpectedCpf: "12345678909",
		},
		{
			Name:        "clean valid CPF format without hyphens and dots",
			InputCpf:    "98765432109",
			ExpectedCpf: "98765432109",
		},
		{
			Name:        "clean invalid CPF format with additional characters",
			InputCpf:    "999-888-777-X",
			ExpectedCpf: "999888777X",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := CleanCpf(test.InputCpf)
			assert.Equal(t, test.ExpectedCpf, result)
		})
	}
}
