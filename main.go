package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/Sensor", func(c *fiber.Ctx) error {
		return c.Send([]byte("wasaaa 2"))
	})

	app.Get("/Sensor/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response := fmt.Sprintf("ID: %s", id)
		return c.Send([]byte([]byte(response)))
	})

	app.Listen("localhost:8000")
}
