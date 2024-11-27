package main

import (
	"fmt"
	"log"
	accountsapiv1 "packages/proto-gen/go/accounts/accountsapi/v1"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/protobuf/proto"
)

func Hello(name string) string {
	result := "Hello " + name
	return result
}

func main() {
	fmt.Println(Hello("accounts-api"))

	person := &accountsapiv1.Person{FirstName: "Eric", LastName: "Zorn", Age: 29}
	b, _ := proto.Marshal(person)
	fmt.Println("bytes", b)

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello world from Accounts API!")
	})

	log.Fatal(app.Listen(":3000"))
}
