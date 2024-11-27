package main

import (
	"fmt"
	"log"

	accountsapiv1 "github.com/ericzorn93/career-cue/proto-gen/go/accounts/accountsapi/v1"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/protobuf/proto"
)

func main() {
	person := &accountsapiv1.Person{FirstName: "Eric", LastName: "Zorn", Age: 29}
	b, _ := proto.Marshal(person)
	fmt.Println("bytes", b)

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello world from Accounts API!")
	})

	log.Fatal(app.Listen(":3000", fiber.ListenConfig{
		EnablePrefork: true,
	}))
}
