package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.io/victor99z/rinha/database"
	"github.io/victor99z/rinha/handlers"
	"github.io/victor99z/rinha/repository"
)

func main() {

	db := database.CreateConnection()
	storage := repository.NewClienteStorage(db)
	handlers := handlers.NewClienteHandlers(storage)

	defer db.Close()

	// Muda o decoder e encoder padrão do Fiber para o go-json
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Post("/clientes/:id/transacoes", handlers.NewTransaction)
	app.Get("/clientes/:id/extrato", handlers.GetBalance)

	log.Fatal(app.Listen(":3000"))

}
