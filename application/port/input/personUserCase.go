package input

import "github.com/petshop-system/petshop-api/application/domain"

type IPersonService interface {
	Create(contextControl domain.ContextControl, customer domain.PersonDomain) (domain.PersonDomain, error)
}
