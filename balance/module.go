package balance

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance/application"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/config/env"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/server/web_server"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/api/handle"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_handle"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/log_init"
	"os"
)

/**
  inject and summary all logic construct service, map application, router, service,...
*/

type (
	Module struct {
		HttpServer web_server.HttpServer
		RouterHttp *Routes
		Service    application.ServiceInterface
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
	m.ResignApi()
}

func (m *Module) Inject() {
	m.HttpServer = m.NewWebServer()
	m.Service = ForwardProviderService()
	m.RouterHttp = m.NewRouter()
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
	m.HttpServer.Get("balance/health-check", m.RouterHttp.apiController.HealthCheck)
	m.HttpServer.Post("balance/request-balance", m.RouterHttp.apiController.HandleOneRequestBalance)
}

func (m *Module) StartWebServer() {
	errApp := m.HttpServer.Listen(":8080")
	if errApp != nil {
		errMApp := fmt.Sprintf("Init app error, message: %s", errApp.Error())
		panic(errMApp)
		os.Exit(0)
	}
}

func (m *Module) Init() {
	m.InitLogs()
	m.InitEnv()
	m.Inject()
	m.InitService()
	m.ResignRoutes()
}

func (m *Module) InitService() {
	m.Service.Init()
}
func (m *Module) InitEnv() {
	linkENVLocal, statusELC := os.LookupEnv("LINK_ENV_LOCAL")
	if statusELC == false {
		panic("missing env LINK_ENV_LOCAL")
	}

	gCf := env.NewGlobalConfig()
	gCf.Init(env.ConfigEnv{
		FilePatchLogInLocal: linkENVLocal,
	})
}

func (m *Module) InitLogs() {
	linkFileLog, statusLLF := os.LookupEnv("LINK_FILE_LOG")
	if statusLLF == false {
		panic("missing env LINK_LOG_FILE")
	}

	linkFolderLog, statusLPL := os.LookupEnv("LINK_FOLDER_LOG")
	if statusLPL == false {
		panic("missing env LINK_FOLDER_LOG")
	}

	log_init.Init(log_init.Log{
		TypeFormat: log_init.TYPE_FORMAT_TEXT,
		TypeOutput: log_init.TYPE_OUTPUT_FILE,
		LinkFile:   linkFileLog,
		LinkFolder: linkFolderLog,
	})
}

func (m *Module) NewRouter() *Routes {
	r := Routes{
		viewController: nil,
		apiController:  handle.NewRequestBalance(m.Service),
	}

	return &r
}

//call if use micro service, one project use one module
//if in monothic, custom  param pass to Init() and run all modul in main.go
func (m *Module) Start() {
	m.StartWebServer()
}

func NewModule() *Module {
	var _ ModuleInterface = (*Module)(nil)
	m := Module{}
	m.Init()
	return &m
}

func NewRouter(viewController any, apiController *handle.RequestBalance) *Routes {
	return &Routes{
		viewController: viewController,
		apiController:  apiController,
	}
}
