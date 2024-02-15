package models

type Transaction struct {
	Id        int    `db:"id"`
	ClienteId string `db:"cliente_id" validate:"required"`
	Valor     int    `db:"valor" validate:"required,gte=0"`
	Tipo      string `db:"tipo" validate:"required,min=1,max=1"`
	Descricao string `db:"descricao" validate:"required,min=1,max=10"`
	Date      string `db:"realizada_em"`
}

type Cliente struct {
	Id     int    `db:"id"`
	Nome   string `db:"nome"`
	Limite int    `db:"limite"`
}

type Saldo struct {
	Id        int `db:"id"`
	ClienteId int `db:"cliente_id"`
	Valor     int `db:"valor"`
}
