package input

import "github.com/petshop-system/petshop-api/application/domain"

type ICustomerService interface {
	Create(contextControl domain.ContextControl, customer domain.CustomerDomain) (domain.CustomerDomain, error)
	ValidateTypePerson(customer domain.CustomerDomain) error
	ValidateCreate(customer domain.CustomerDomain) error
}
