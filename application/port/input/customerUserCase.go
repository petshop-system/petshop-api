package input

import "github.com/petshop-system/petshop-api/application/domain"

type ICustomerService interface {
	Create(contextControl domain.ContextControl, customer domain.Customer) (domain.Customer, error)
}
