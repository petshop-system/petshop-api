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
			InputCpf:      "12345678909",
			ExpectedError: "",
		},
		{
			Name:          "invalid CPF format",
			InputCpf:      "123-456-789",
			ExpectedError: ErrorInvalidCPFLength,
		},
		{
			Name:          "invalid CPF with all digits equal",
			InputCpf:      "111.111.111-11",
			ExpectedError: ErrorAllDigitsEqualCPF,
		},
		{
			Name:          "valid CPF with extra characters",
			InputCpf:      "123.456.789-09-X",
			ExpectedError: ErrorInvalidCPFLength,
		},
		{
			Name:          "invalid CPF with incorrect first verification digit",
			InputCpf:      "076.164.346-16",
			ExpectedError: ErrorFirstVerificationCPF,
		},
		{
			Name:          "invalid CPF with incorrect second verification digits",
			InputCpf:      "529.982.247-27",
			ExpectedError: ErrorSecondVerificationCPF,
		},
		{
			Name:          "invalid CPF with incorrect character conversion",
			InputCpf:      "12A.345.678-90",
			ExpectedError: ErrorIncorrectCharacterConversionCPF,
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

func TestRemoveNonNumericCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123-456-789", "123456789"},
		{"1.2.3", "123"},
	}

	for _, test := range tests {
		result := RemoveNonNumericCharacters(test.input)
		if result != test.expected {
			t.Errorf("RemoveNonNumericCharacters(%s) = %s, esperado %s", test.input, result, test.expected)
		}
	}
}
