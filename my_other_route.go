package main

import "github.com/gofiber/fiber/v2"

func UserRoutes() *fiber.App {
	micro := fiber.New()

	micro.Get("/doe", func(c *fiber.Ctx) error {
		return c.SendString("DOEEEE")
	})

	return micro
}

func MountOtherRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/list", func(c *fiber.Ctx) error {
		return c.SendString("grouped route /api/v1/list")
	})
}
