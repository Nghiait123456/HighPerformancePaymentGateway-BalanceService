package balance

import (
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/server/web_server"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/api/handle"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_handle"
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
		Inject()
		Init()
	}

	Routes struct {
		viewController any // project don't have view, only rest api
		apiController  *handle.RequestBalance
	}
)

func (m *Module) ResignRoutes() {

}

func (m *Module) Inject() {
	m.HttpServer = m.NewWebServer()

}

func (m Module) NewWebServer() web_server.HttpServer {
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: error_handle.CustomMessageError,
	})

	//config error handle
	eH := error_handle.ErrorHandle{
		App: app,
	}
	eH.Init()

	// config alert panic notify
	alertAc := error_handle.AlertAc{
		Dsn:              "https://4ea29cdaaa424266a113571ac407c88a@o1406092.ingest.sentry.io/6739879",
		Repanic:          true,
		Debug:            true,
		AttachStacktrace: true,
	}
	alertPanic := error_handle.NewPanicHandle(&alertAc, app)
	alertPanic.Init()

	return app
}

func (m *Module) ResignApi() {
	m.HttpServer.Post("/request-balance", m.RouterHttp.apiController.HandleOneRequestBalance)
}

func (m *Module) Init() {
	m.Inject()
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
