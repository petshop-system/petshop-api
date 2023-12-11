package input

import "github.com/petshop-system/petshop-api/application/domain"

type IPersonService interface {
	Create(contextControl domain.ContextControl, customer domain.PersonDomain) (domain.PersonDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.PersonDomain, bool, error)
}
