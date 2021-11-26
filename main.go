package main

import (
	"api-cache-proxy/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"os"
	"strings"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	if strings.EqualFold("TRUE", os.Getenv("SHOW_REQUEST_LOG")) {
		app.Use(logger.New())
	}
	app.Use(services.NewCaching())
	app.Get("/*", func(c *fiber.Ctx) error {
		url := os.Getenv("TARGET_HOST") + string(c.Request().RequestURI())
		if err := proxy.Do(c, url); err != nil {
			return err
		}
		// Remove Server header from response
		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})
	app.Listen(":3000")
}
