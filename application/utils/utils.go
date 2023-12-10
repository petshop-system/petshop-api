package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ErrorInvalidCPFLength                = "invalid CPF length error"
	ErrorAllDigitsEqualCPF               = "invalid CPF because all digits are equal"
	ErrorFirstVerificationCPF            = "error in the first verification of the CPF"
	ErrorSecondVerificationCPF           = "error in the second verification of the CPF"
	ErrorIncorrectCharacterConversionCPF = "invalid CPF with incorrect character conversion"
)

func ValidateCpf(cpf string) error {
	cleanedCpf := RemoveNonNumericCharacters(cpf)

	if len(cleanedCpf) != 11 {
		return fmt.Errorf(ErrorInvalidCPFLength)
	}

	characters := strings.Split(cleanedCpf, "")
	firstVerification, _ := strconv.Atoi(characters[9])
	secondVerification, _ := strconv.Atoi(characters[10])

	allDigitsEqual := strings.Count(cleanedCpf, string(cleanedCpf[0])) == 11

	if allDigitsEqual {
		return fmt.Errorf(ErrorAllDigitsEqualCPF)
	}

	var status1, status2 int

	for i := 0; i < 9; i++ {
		num, err := strconv.Atoi(characters[i])
		if err != nil {
			return fmt.Errorf(ErrorIncorrectCharacterConversionCPF)
		}
		status1 += num * (10 - i)
	}
	test1 := (status1 * 10) % 11

	if test1 == 10 {
		test1 = 0
	}

	if test1 != firstVerification {
		return fmt.Errorf(ErrorFirstVerificationCPF)
	}

	for i := 0; i < 10; i++ {
		num, err := strconv.Atoi(characters[i])
		if err != nil {
			return fmt.Errorf(ErrorIncorrectCharacterConversionCPF)
		}

		status2 += num * (11 - i)
	}
	test2 := (status2 * 10) % 11

	if test2 != secondVerification {
		return fmt.Errorf(ErrorSecondVerificationCPF)
	}
	return nil
}

func RemoveNonNumericCharacters(documentNumber string) string {
	documentNumber = strings.ReplaceAll(documentNumber, ".", "")
	documentNumber = strings.ReplaceAll(documentNumber, "-", "")
	documentNumber = strings.ReplaceAll(documentNumber, "/", "")
	return documentNumber
}

/*func ValidateCnpj(cnpj string) error {
	return nil
}*/
