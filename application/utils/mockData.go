package utils

import "github.com/petshop-system/petshop-api/application/domain"

const (
	MockAddressStreet       = "Rua Conde de Bonfim"
	MockAddressNumber       = "123"
	MockAddressComplement   = "303"
	MockAddressNeighborhood = "Tijuca"
	MockAddressZipCode      = "20520-050"
	MockAddressCity         = "Rio de Janeiro"
	MockAddressState        = "RJ"
	MockAddressCountry      = "Brasil"
)

func GetMockAddress() domain.AddressDomain {
	return domain.AddressDomain{
		Street:       MockAddressStreet,
		Number:       MockAddressNumber,
		Complement:   MockAddressComplement,
		Neighborhood: MockAddressNeighborhood,
		ZipCode:      MockAddressZipCode,
		City:         MockAddressCity,
		State:        MockAddressState,
		Country:      MockAddressCountry,
	}
}
