package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
)

func Hello(name string) string {
	result := "Hello " + name
	return result
}

func main() {
	fmt.Println(Hello("inbound-webhooks-api"))

	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello world from Inbound Webhooks API!")
	})

	log.Fatal(app.Listen(":3000"))
}
