package error_handle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type (
	ErrorHandle struct {
		App *fiber.App
	}

	ErrorHandleInterface interface {
		Init()
	}
)

func (e *ErrorHandle) Init() {
	e.App.Use(recover.New())
}

/**
if want custom, pass function to config

App := fiber.New(fiber.Config{
// Override default error handler
ErrorHandler: CustomMessageError})
*/
func CustomMessageError(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Send custom error page
	err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(fiber.StatusInternalServerError).SendString("Something went wrong, Internal Server Error")
	}

	// Return from handler
	return nil
}
