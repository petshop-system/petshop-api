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

const (
	ErrorInvalidMobilePhoneLength   = "invalid Mobile Phone length error"
	ErrorInvalidLandLinePhoneLength = "invalid Land Line Phone length error"
	ErrorAreaCodeVerification       = "invalid area code"
)

const (
	LandLinePhone = "landline_phone"
	MobilePhone   = "mobile_phone"
)

func ValidatePhoneNumber(phoneType, phoneNumber string) error {

	clearPhone := RemoveNonAlphaNumericCharacters(phoneNumber)
	verification := func(phoneLen int, phoneTypeVerification, errorMessageVerification string) error {
		if len(clearPhone) != phoneLen {
			return fmt.Errorf(errorMessageVerification)
		}
		return nil
	}

	switch phoneType {
	case LandLinePhone:
		landLinePhoneLen := 8
		if err := verification(landLinePhoneLen, phoneType, ErrorInvalidLandLinePhoneLength); err != nil {
			return err
		}
	case MobilePhone:
		mobilePhoneLen := 9
		if err := verification(mobilePhoneLen, phoneType, ErrorInvalidMobilePhoneLength); err != nil {
			return err
		}
	}
	return nil
}

func ValidateCodeAreaNumber(areaCode string) error {

	dddList := map[string]bool{
		// North Region
		"68": true, "96": true, "92": true, "97": true, "91": true, "93": true, "94": true, "69": true, "95": true, "63": true,
		// Northeast Region
		"82": true, "71": true, "73": true, "74": true, "75": true, "77": true, "85": true, "88": true, "98": true, "99": true, "83": true, "81": true, "87": true, "86": true, "89": true, "84": true, "79": true,
		// Midwest Region
		"61": true, "62": true, "64": true, "65": true, "66": true,
		// Southeast Region
		"27": true, "28": true, "31": true, "32": true, "33": true, "34": true, "35": true, "37": true, "38": true, "21": true, "22": true, "24": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		// South Region
		"41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "51": true, "53": true, "54": true, "55": true,
	}

	clearAreaCode := RemoveNonAlphaNumericCharacters(areaCode)

	if dddList[clearAreaCode] {
		return nil
	}
	return fmt.Errorf(ErrorAreaCodeVerification)
}

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
