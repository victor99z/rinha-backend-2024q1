package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.io/victor99z/rinha/database"
	"github.io/victor99z/rinha/handlers"
)

func DebugHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func main() {

	db, _ := database.CreateConnection()

	// Cria uma nova instância do ClienteStorage
	storage := handlers.NewClienteStorage(db)

	// Muda o decoder e encoder padrão do Fiber para o go-json
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Post("/clientes/:id/transacoes", storage.NewTransaction)
	app.Get("/clientes/:id/extrato", storage.GetBalance)

	log.Fatal(app.Listen(":3000"))
}
