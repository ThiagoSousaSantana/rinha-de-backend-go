package models

type Cliente struct {
	ID     int `json:"id"`
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

type Transacao struct {
	ID        int    `json:"id"`
	ClienteID int    `json:"cliente_id"`
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
	Data      string `json:"data"`
}
