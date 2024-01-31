package input

import "github.com/petshop-system/petshop-api/application/domain"

type IPhoneService interface {
	Create(contextControl domain.ContextControl, phone domain.PhoneDomain) (domain.PhoneDomain, error)
	GetByID(contextControl domain.ContextControl, ID int64) (domain.PhoneDomain, bool, error)
}
