package repositories

import (
	"github.com/thiagosousasantana/rinha-go/db"
	"github.com/thiagosousasantana/rinha-go/models"
)

func InsertTransacao(transacao models.Transacao) error {

	sql := `INSERT INTO transacoes (cliente_id, valor, tipo, data, descricao) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.CONN.Exec(sql, transacao.ClienteID, transacao.Valor, transacao.Tipo, transacao.Data, transacao.Descricao)

	if err != nil {
		return err
	}

	return nil
}

func FindAllByClienteId(clienteId int8) (transacoes []models.Transacao, err error) {

	sql := `SELECT id, cliente_id, valor, tipo, data, descricao FROM transacoes WHERE cliente_id = $1`

	rows, err := db.CONN.Query(sql, clienteId)
	if err != nil {
		return
	}
	defer rows.Close()
	transacoes = []models.Transacao{}

	for rows.Next() {
		var transacao models.Transacao
		err = rows.Scan(&transacao.ID, &transacao.ClienteID, &transacao.Valor, &transacao.Tipo, &transacao.Data, &transacao.Descricao)
		if err != nil {
			return
		}

		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}

func FindLast10ByClienteId(clienteId int8) (transacoes []models.Transacao, err error) {

	sql := `SELECT id, cliente_id, valor, tipo, data, descricao FROM transacoes WHERE cliente_id = $1 ORDER BY id DESC LIMIT 10`

	rows, err := db.CONN.Query(sql, clienteId)
	if err != nil {
		return
	}
	defer rows.Close()
	transacoes = []models.Transacao{}

	for rows.Next() {
		var transacao models.Transacao
		err = rows.Scan(&transacao.ID, &transacao.ClienteID, &transacao.Valor, &transacao.Tipo, &transacao.Data, &transacao.Descricao)
		if err != nil {
			return
		}

		transacoes = append(transacoes, transacao)
	}
	return transacoes, nil
}
