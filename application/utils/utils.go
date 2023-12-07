package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	ErrorInvalidCpfFormat         = "invalid CPF format"
	ErrorAllDigitsEqual           = "invalid CPF because all digits are equal"
	ErrorFirstVerification        = "error in the first verification of the CPF"
	ErrorSecondVerification       = "error in the second verification of the CPF"
	ErrorConvertingCharacterToNum = "error converting character to number: %v"
)

func ValidateCpf(cpf string) error {
	if !isValidCpfFormat(cpf) {
		return fmt.Errorf(ErrorInvalidCpfFormat)
	}

	cleanedCpf := CleanCpf(cpf)

	characters := strings.Split(cleanedCpf, "")
	firstVerification, _ := strconv.Atoi(characters[9])
	secondVerification, _ := strconv.Atoi(characters[10])

	allDigitsEqual := strings.Count(cleanedCpf, string(cleanedCpf[0])) == 11

	if allDigitsEqual {
		return fmt.Errorf(ErrorAllDigitsEqual)
	}

	var status1, status2 int

	for i := 0; i < 9; i++ {
		num, err := strconv.Atoi(characters[i])
		if err != nil {
			return fmt.Errorf(ErrorConvertingCharacterToNum, err)
		}
		status1 += num * (10 - i)
	}
	test1 := (status1 * 10) % 11

	if test1 == 10 {
		test1 = 0
	}

	if test1 != firstVerification {
		return fmt.Errorf(ErrorFirstVerification)
	}

	for i := 0; i < 10; i++ {
		num, err := strconv.Atoi(characters[i])
		if err != nil {
			return fmt.Errorf(ErrorConvertingCharacterToNum, err)
		}

		status2 += num * (11 - i)
	}
	test2 := (status2 * 10) % 11

	if test2 != secondVerification {
		return fmt.Errorf(ErrorSecondVerification)
	}
	return nil
}

func isValidCpfFormat(cpf string) bool {
	patternCpf := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$|^\d{11}$`)
	return patternCpf.MatchString(cpf)
}

func CleanCpf(cpf string) string {
	return strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")
}

/*func ValidateCnpj(cnpj string) error {
	return nil
}*/
