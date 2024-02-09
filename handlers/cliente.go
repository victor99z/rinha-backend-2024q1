package handlers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.io/victor99z/rinha/models"
)

type ClienteStorage struct {
	Conn *sqlx.DB
}

func NewClienteStorage(conn *sqlx.DB) *ClienteStorage {
	return &ClienteStorage{Conn: conn}
}

type ExtratoFormatado struct {
	Saldo models.Saldo         `json:"saldo"`
	Tx    []models.Transaction `json:"ultimas_transacoes"`
}

type (
	TxDTO struct {
		Valor     int64  `json:"valor"`
		Tipo      string `json:"tipo"`
		Descricao string `json:"descricao"`
		Date      string `json:"realizada_em"`
	}
	SaldoDTO struct {
		Total       int64  `json:"total"`
		DataExtrato string `json:"data_extrato"`
		Limite      int64  `json:"limite"`
	}
)

func (s *ClienteStorage) GetBalance(c *fiber.Ctx) error {

	saldoCliente := new([]models.Saldo)
	ultimasTx := new([]models.Transaction)

	errSaldo := s.Conn.Select(saldoCliente, "SELECT * FROM saldos WHERE cliente_id = $1", c.Params("id"))

	if errSaldo != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(errSaldo.Error())
	}

	firstSaldoClinte := (*saldoCliente)[0]

	if firstSaldoClinte.Id == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	errTx := s.Conn.Select(ultimasTx, "SELECT * FROM transacoes WHERE cliente_id = $1 ORDER BY realizada_em DESC LIMIT 10", c.Params("id"))

	if errTx != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	ExtratoFormatado := ExtratoFormatado{
		Saldo: firstSaldoClinte,
		Tx:    *ultimasTx,
	}

	return c.Status(fiber.StatusOK).JSON(ExtratoFormatado)
}

func (s *ClienteStorage) NewTransaction(c *fiber.Ctx) error {

	newTx := new(models.Transaction)

	if err := c.BodyParser(&newTx); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	validate := validator.New()
	errValidation := validate.Struct(newTx)

	if errValidation != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	newTx.ClienteId = (c.Params("id"))

	// s.Conn.MustExec("UPDATE saldo SET valor = valor + $1 WHERE cliente_id = $2", newTx.Valor, newTx.ClienteId)

	res := s.Conn.MustExec("INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES ($1, $2, $3, $4)", newTx.ClienteId, newTx.Valor, newTx.Tipo, newTx.Descricao)

	if v, _ := res.RowsAffected(); v == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusCreated)
}
