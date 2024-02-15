package repository

import (
	"errors"
	"math"

	_ "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"github.io/victor99z/rinha/models"
)

type ClienteStorage struct {
	Conn *sqlx.DB
}

func NewClienteStorage(conn *sqlx.DB) *ClienteStorage {
	return &ClienteStorage{Conn: conn}
}

type SaldoDTO struct {
	Total       int    `json:"total" db:"total"`
	DataExtrato string `json:"data_extrato" db:"data_extrato"`
	Limite      int    `json:"limite" db:"limite"`
}
type TxDTO struct {
	Valor       int    `json:"valor" db:"valor"`
	Tipo        string `json:"tipo" db:"tipo"`
	Descricao   string `json:"descricao" db:"descricao"`
	RealizadaEm string `json:"realizada_em" db:"realizada_em"`
}

type Extrato struct {
	Saldo     SaldoDTO `json:"saldo"`
	UltimasTx []TxDTO  `json:"ultimas_transacoes"`
}

func (s *ClienteStorage) GetExtrato(clienteId string) (Extrato, error) {

	var resTx []TxDTO
	var resSaldo SaldoDTO

	err := s.Conn.Get(&resSaldo, "SELECT valor as total, limite, now() as data_extrato FROM saldos as s, clientes as c WHERE c.id = $1 and s.cliente_id = $1", clienteId)

	if err != nil || resSaldo.Limite == 0 {
		return Extrato{}, errors.New("cliente not found")
	}

	err = s.Conn.Select(&resTx, "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id = $1 ORDER BY realizada_em DESC LIMIT 10", clienteId)

	if err != nil {
		return Extrato{}, err
	}

	return Extrato{
		Saldo:     resSaldo,
		UltimasTx: resTx,
	}, nil
}

type ReturnTx struct {
	Saldo  int `json:"saldo" db:"valor"`
	Limite int `json:"limite" db:"limite"`
}

func (s *ClienteStorage) CreateTx(transaction *models.Transaction) (ReturnTx, error) {

	response := ReturnTx{}

	tx, _ := s.Conn.Beginx()

	infoCliente := tx.Get(&response, "SELECT valor, limite FROM saldos, clientes WHERE saldos.cliente_id = $1 and clientes.id = $1", transaction.ClienteId)

	if infoCliente != nil {
		log.Error(infoCliente.Error())
		tx.Rollback()
		return ReturnTx{}, errors.New("cliente not found")
	}

	if transaction.Tipo == "d" && float64(response.Limite) < math.Abs(float64(response.Saldo-transaction.Valor)) {
		tx.Rollback()
		return ReturnTx{}, errors.New("saldo invalido")
	}

	tx.Exec("INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES ($1, $2, $3, $4)", transaction.ClienteId, transaction.Valor, transaction.Tipo, transaction.Descricao)

	if transaction.Tipo == "c" {
		tx.Exec("UPDATE saldos SET valor = valor + $1 WHERE cliente_id = $2", transaction.Valor, transaction.ClienteId)
	} else {
		tx.Exec("UPDATE saldos SET valor = valor - $1 WHERE cliente_id = $2", transaction.Valor, transaction.ClienteId)
	}

	tx.Commit()

	tx.Get(&response, "SELECT valor, limite FROM saldos WHERE cliente_id = $1", transaction.ClienteId)

	return response, nil
}
