package error_handle

import (
	"fmt"
	sentryfiber "github.com/aldy505/sentry-fiber"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type (
	PanicHandle struct {
		alert *AlertAc
		app   *fiber.App
	}

	AlertAc = Sentry
	Sentry  struct {
		Dsn              string
		Debug            bool
		AttachStacktrace bool
		Repanic          bool
		WaitForDelivery  bool
		Timeout          time.Duration
	}

	PanicHandleInterface interface {
		Init()
	}
)

func (p *PanicHandle) Init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              p.alert.Dsn,
		Debug:            p.alert.Debug,
		AttachStacktrace: p.alert.AttachStacktrace,
	})
	if err != nil {
		errMessage := fmt.Sprintf("init Sentry error")
		log.WithFields(log.Fields{
			"errMessage": errMessage,
		}).Error(errMessage)

		panic(errMessage)
		os.Exit(0)
	}

	//resign middleware capture panic
	p.app.Use(sentryfiber.New(sentryfiber.Options{
		Repanic:         p.alert.Repanic,
		WaitForDelivery: p.alert.WaitForDelivery,
		Timeout:         p.alert.Timeout,
	}))
}

func NewPanicHandle(alert *AlertAc, app *fiber.App) PanicHandleInterface {
	return &PanicHandle{
		alert: alert,
		app:   app,
	}
}
