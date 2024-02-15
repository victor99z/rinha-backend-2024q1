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

type TxDTO struct {
	ClienteId string `db:"cliente_id"`
	Valor     int    `db:"valor" validate:"required,gt=0"`
	Tipo      string `db:"tipo" validate:"required,len=1"`
	Descricao string `db:"descricao" validate:"required,min=1,max=10"`
}

func (s *ClienteHandlers) NewTransaction(c *fiber.Ctx) error {

	newTx := new(TxDTO)

	if err := c.BodyParser(&newTx); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if validate := validator.New(); validate.Struct(newTx) != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if newTx.Tipo != "c" && newTx.Tipo != "d" {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	newTx.ClienteId = (c.Params("id"))

	res, err := s.Store.CreateTx(&models.Transaction{
		ClienteId: newTx.ClienteId,
		Valor:     newTx.Valor,
		Tipo:      newTx.Tipo,
		Descricao: newTx.Descricao,
	})

	if err != nil {
		if err.Error() == "saldo invalido" {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if err.Error() == "cliente not found" {
			return c.SendStatus(fiber.StatusNotFound)
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
