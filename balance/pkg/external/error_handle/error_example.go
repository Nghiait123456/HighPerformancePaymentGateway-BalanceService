package error_handle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/log_init"
	"os"
)

// handle panic ==> alert message, custom message 500 internal server
func example() {
	log_init.Init(log_init.Log{
		TypeFormat: log_init.TYPE_FORMAT_TEXT,
		TypeOutput: log_init.TYPE_OUTPUT_FILE,
		LinkFile:   "balance/infrastructure/log/log_file/log.log",
	})
	//
	//gCf := TestEnvGlobal()
	//TestAutoChangeSecret(gCf.AuthInternalServiceConfig())

	// Create a new fiber instance with custom config
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: CustomMessageError,
	})

	//config error handle
	eH := ErrorHandle{
		App: app,
	}
	eH.Init()

	// config alert panic notify
	alertAc := AlertAc{
		Dsn:              "https://4ea29cdaaa424266a113571ac407c88a@o1406092.ingest.sentry.io/6739879",
		Repanic:          true,
		Debug:            true,
		AttachStacktrace: true,
	}
	alertPanic := NewPanicHandle(&alertAc, app)
	alertPanic.Init()

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("panic0999hbhbhb" + "333333333333333333333")
	})

	errApp := app.Listen(":3000")
	if errApp != nil {
		errMApp := fmt.Sprintf("Init app error, message: %s", errApp.Error())
		panic(errMApp)
		os.Exit(0)
	}
}
