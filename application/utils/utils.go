package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	ErrorInvalidCPFLength      = "invalid CPF length error"
	ErrorAllDigitsEqualCPF     = "invalid CPF because all digits are equal"
	ErrorFirstVerificationCPF  = "error in the first verification of the CPF"
	ErrorSecondVerificationCPF = "error in the second verification of the CPF"

	ErrorInvalidCNPJLength      = "invalid CNPJ length error"
	ErrorAllDigitsEqualCNPJ     = "invalid CNPJ because all digits are equal"
	ErrorFirstVerificationCNPJ  = "error in the first verification of the CNPJ"
	ErrorSecondVerificationCNPJ = "error in the second verification of the CNPJ"
)

func ValidateCpf(cpf string) error {
	cleanedCpf := RemoveNonAlphaNumericCharacters(cpf)

	cpfLen := 11
	if len(cleanedCpf) != cpfLen {
		return fmt.Errorf(ErrorInvalidCPFLength)
	}

	if allDigitsEqual := strings.Count(cleanedCpf, string(cleanedCpf[0])) == cpfLen; allDigitsEqual {
		return fmt.Errorf(ErrorAllDigitsEqualCPF)
	}

	characters := strings.Split(cleanedCpf, "")
	beforeLastPosition := len(cleanedCpf) - 2
	beforeLastValue, _ := strconv.Atoi(characters[beforeLastPosition])
	lastPosition := len(cleanedCpf) - 1
	lastValue, _ := strconv.Atoi(characters[lastPosition])

	/*
		positionVerification: the number of before last char position or last char position
		lenCharacters: size of characters in the cpf
		charVerification: the value of before last char position or last char position
		errorMessageVerification: message to error char verification
	*/
	verification := func(positionVerification, lenCharacters, charVerification int, errorMessageVerification string) error {

		var status int
		constValidation := 10

		for i := 0; i < positionVerification; i++ {
			num, err := strconv.Atoi(characters[i])
			if err != nil {
				return fmt.Errorf(ErrorInvalidCPFLength)
			}
			status += num * ((positionVerification + 1) - i)
		}

		checkCharacter := (status * constValidation) % lenCharacters

		if checkCharacter == constValidation {
			checkCharacter = 0
		}

		if checkCharacter != charVerification {
			return fmt.Errorf(errorMessageVerification)
		}

		return nil
	}

	if err := verification(beforeLastPosition, cpfLen, beforeLastValue, ErrorFirstVerificationCPF); err != nil {
		return err
	}

	if err := verification(lastPosition, cpfLen, lastValue, ErrorSecondVerificationCPF); err != nil {
		return err
	}

	return nil
}

func ValidateCnpj(cnpj string) error {
	cleanedCnpj := RemoveNonAlphaNumericCharacters(cnpj)

	cnpjLen := 14
	if len(cleanedCnpj) != cnpjLen {
		return fmt.Errorf(ErrorInvalidCNPJLength)
	}

	if allDigitsEqual := strings.Count(cleanedCnpj, string(cleanedCnpj[0])) == cnpjLen; allDigitsEqual {
		return fmt.Errorf(ErrorAllDigitsEqualCNPJ)
	}

	characters := strings.Split(cleanedCnpj, "")
	beforeLastPosition := len(cleanedCnpj) - 2
	beforeLastValue, _ := strconv.Atoi(characters[beforeLastPosition])
	lastPosition := len(cleanedCnpj) - 1
	lastValue, _ := strconv.Atoi(characters[lastPosition])

	verification := func(positionVerification int, multipliers []int, charVerification int, errorMessageVerification string) error {
		var status int

		for i := 0; i < positionVerification; i++ {
			num, err := strconv.Atoi(characters[i])
			if err != nil {
				return fmt.Errorf(ErrorInvalidCNPJLength)
			}
			status += num * multipliers[i]
		}

		const cnpjCalculationFactor = 11
		minValueVerification := 2
		maxValueVerification := 10

		checkCharacter := status % cnpjCalculationFactor

		if (checkCharacter == 0 || checkCharacter == 1) && charVerification != 0 { //Calculation Rule for CNPJ. If the remainder of the division is 0 or 1, the penultimate digit should be 0.
			return fmt.Errorf(errorMessageVerification)
		}

		if (checkCharacter >= minValueVerification && checkCharacter <= maxValueVerification) && (charVerification != cnpjCalculationFactor-checkCharacter) { //from min to max: Calculation Rule for CNPJ
			return fmt.Errorf(errorMessageVerification)
		}
		return nil
	}

	cnpjMultipliersFirst := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2} // Calculation Rule for CNPJ. Array for the first step of CNPJ calculation.
	if err := verification(beforeLastPosition, cnpjMultipliersFirst, beforeLastValue, ErrorFirstVerificationCNPJ); err != nil {
		return err
	}

	cnpjMultipliersSecond := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2} //Calculation Rule for CNPJ. Array for the second step of CNPJ calculation.
	if err := verification(lastPosition, cnpjMultipliersSecond, lastValue, ErrorSecondVerificationCNPJ); err != nil {
		return err
	}
	return nil
}

func RemoveNonAlphaNumericCharacters(documentNumber string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(documentNumber, "")
}
