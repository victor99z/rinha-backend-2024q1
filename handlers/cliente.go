package handlers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.io/victor99z/rinha/models"
	"github.io/victor99z/rinha/repository"
)

type ClienteHandlers struct {
	Store *repository.ClienteStorage
}

func NewClienteHandlers(store *repository.ClienteStorage) *ClienteHandlers {
	return &ClienteHandlers{Store: store}
}

func (s *ClienteHandlers) GetBalance(c *fiber.Ctx) error {

	res, err := s.Store.GetExtrato(c.Params("id"))

	if err != nil {
		if err.Error() == "cliente not found" {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *ClienteHandlers) NewTransaction(c *fiber.Ctx) error {

	newTx := new(models.Transaction)

	if err := c.BodyParser(&newTx); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if validate := validator.New(); validate.Struct(newTx) != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	newTx.ClienteId = (c.Params("id"))

	res, err := s.Store.CreateTx(newTx)

	if err != nil {
		if err.Error() == "saldo invalido" {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if err.Error() == "cliente not found" {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.SendStatus(fiber.StatusUnprocessableEntity)

	}

	return c.Status(fiber.StatusOK).JSON(res)
}
