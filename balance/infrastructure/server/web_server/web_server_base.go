package web_server

import "github.com/gofiber/fiber/v2"

type (
	ContextBase = *fiber.Ctx
	MapBase     = fiber.Map
	HttpServer  = *fiber.App
)
