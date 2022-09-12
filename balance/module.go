package balance

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/server/web_server"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/api/handle"
)

/**
  inject and summary all logic construct service, map application, router, service,...
*/

type (
	Module struct {
		HttpServer web_server.HttpServer
		RouterHttp *Routes
		// todo other config
		// global value
	}

	ModuleInterface interface {
		ResignRoutes()
		ResignApi()
		Init()
	}

	Routes struct {
		viewController any // project don't have view, only rest api
		apiController  *handle.RequestBalance
	}
)

func (m *Module) ResignRoutes() {

}

func (m *Module) ResignApi() {
	m.HttpServer.Post("/request-balance", m.RouterHttp.apiController.HandleOneRequestBalance)
}

func (m *Module) Init() {
	m.ResignRoutes()
	m.ResignApi()
}

func NewModule(httpServer web_server.HttpServer, routerHttp *Routes) *Module {
	var _ ModuleInterface = (*Module)(nil)
	return &Module{
		HttpServer: httpServer,
		RouterHttp: routerHttp,
	}
}
