package main

import (
	"apps/services/accounts-api/internal/random"
	"fmt"
	accountsapiv1 "libs/proto-gen/go/accounts/accountsapi/v1"
	"log"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/protobuf/proto"
)

func main() {
	person := &accountsapiv1.Person{FirstName: "Eric", LastName: "Zorn", Age: 29}
	b, _ := proto.Marshal(person)
	fmt.Println("bytes", string(b))

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		random.RandomMessage("Hello world from Accounts API!")
		return c.SendString("Hello world from Accounts API!")
	})

	log.Fatal(app.Listen(":3000"))
}
