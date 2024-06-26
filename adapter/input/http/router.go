package http

import (
	"github.com/go-chi/chi"
	"github.com/petshop-system/petshop-api/adapter/input/http/handler"
	"go.uber.org/zap"
)

type Router struct {
	ContextPath string
	chiRouter   chi.Router
	LoggerSugar *zap.SugaredLogger
}

func GetNewRouter(loggerSugar *zap.SugaredLogger) Router {
	router := chi.NewRouter()
	return Router{
		chiRouter:   router,
		LoggerSugar: loggerSugar,
	}
}

func (router Router) GetChiRouter() chi.Router {
	return router.chiRouter
}

func (router Router) AddGroupHandlerHealthCheck(ah *handler.Generic) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/health-check", func(r chi.Router) {
			r.Get("/", ah.HealthCheck)
		})
	}
}

func (router Router) AddGroupHandlerCustomer(ah *handler.Customer) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/customer", func(r chi.Router) {
			r.Post("/validate-create", ah.ValidateCreate)
			r.Post("/create", ah.Create)
		})
	}
}

func (router Router) AddGroupHandlerAddress(ah *handler.Address) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/address", func(r chi.Router) {
			r.Post("/create", ah.Create)
			r.Get("/search/{id}", ah.GetByID)
		})
	}
}

func (router Router) AddGroupHandlerPhone(ah *handler.Phone) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/phone", func(r chi.Router) {
			r.Post("/create", ah.Create)
			r.Get("/search/{id}", ah.GetByID)
		})
	}
}
