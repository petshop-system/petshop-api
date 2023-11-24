package input

import "github.com/petshop-system/petshop-api/application/domain"

type ICustomerService interface {
	Create(contextControl domain.ContextControl, customer domain.ClienteDomain) (domain.ClienteDomain, error)
}
