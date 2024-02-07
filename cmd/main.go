package main

import (
	"fmt"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Id  int    `json:"identifier"`
	Msg string `json:"message"`
}

type Transaction struct {
	valor     int64
	tipo      string
	descricao string
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	newUser := User{
		Id:  1,
		Msg: "Hello, World!",
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(newUser)
	})

	app.Post("/clientes/:id/transacoes", func(c *fiber.Ctx) error {

		newTx := new(Transaction)

		if err := c.BodyParser(newTx); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		fmt.Println(c.Body())

		return c.Status(fiber.StatusOK).JSON(map[string]int64{"limite": 10000, "saldo": -9098})
	})

	log.Fatal(app.Listen(":3000"))
}
