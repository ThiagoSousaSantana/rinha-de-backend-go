package repositories

import (
	"database/sql"
	"log"

	"github.com/thiagosousasantana/rinha-go/db"
	"github.com/thiagosousasantana/rinha-go/models"
)

func UpdateSaldo(id int8, saldo int, tx *sql.Tx) (int64, error) {
	sql := `UPDATE clientes SET saldo = $1 WHERE id = $2`

	res, err := tx.Exec(sql, saldo, id)
	if err != nil {
		log.Printf("Error updating saldo. %v", err)
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Printf("Rollback error: %v", rbErr)
		}
		return 0, err
	}

	err = tx.Commit()

	if err != nil {
		log.Printf("Error commiting transaction. %v", err)
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Printf("Rollback error: %v", rbErr)
		}
		return 0, err
	}

	return res.RowsAffected()
}

func FindById(id int) (cliente models.Cliente, err error) {

	sql := `SELECT id, limite, saldo FROM clientes WHERE id = $1`

	err = db.CONN.QueryRow(sql, id).Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)

	return cliente, err
}

func FindByIdLocking(id int, tx *sql.Tx) (cliente models.Cliente, err error) {

	sql := `SELECT id, limite, saldo FROM clientes WHERE id = $1 FOR UPDATE`

	err = tx.QueryRow(sql, id).Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)

	return cliente, err
}
