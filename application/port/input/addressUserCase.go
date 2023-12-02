package input

import "github.com/petshop-system/petshop-api/application/domain"

type IAddressService interface {
	Create(contextControl domain.ContextControl, customer domain.AddressDomain) (domain.AddressDomain, error)
}
