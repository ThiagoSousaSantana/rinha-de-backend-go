package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/thiagosousasantana/rinha-go/config"
)

var CONN *sql.DB

func OpenConnection() {
	conf := config.GetDB()
	var err error

	sc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DataBase)

	CONN, err = sql.Open("postgres", sc)

	if err != nil {
		panic(err)
	}

	CONN.SetMaxOpenConns(50)
	CONN.SetMaxIdleConns(20)

	err = CONN.Ping()

	if err != nil {
		panic(err)
	}
}

func OpenTransaction(ctx context.Context) (*sql.Tx, error) {
	return CONN.BeginTx(ctx, nil)
}
