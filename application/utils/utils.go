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
				return fmt.Errorf(ErrorIncorrectCharacterConversionCPF)
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

func RemoveNonNumericCharacters(documentNumber string) string {
	documentNumber = strings.ReplaceAll(documentNumber, ".", "")
	documentNumber = strings.ReplaceAll(documentNumber, "-", "")
	documentNumber = strings.ReplaceAll(documentNumber, "/", "")
	return documentNumber
}

/*func ValidateCnpj(cnpj string) error {
	return nil
}*/
