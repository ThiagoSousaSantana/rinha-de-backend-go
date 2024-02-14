package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thiagosousasantana/rinha-go/config"
	"github.com/thiagosousasantana/rinha-go/db"
	"github.com/thiagosousasantana/rinha-go/handlers"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	db.OpenConnection()

	r := chi.NewRouter()
	r.Post("/clientes/{id}/transacoes", handlers.CreateTransacao)
	r.Get("/clientes/{id}/extrato", handlers.GetExtrato)

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.GetServerPort()), r)

	if err != nil {
		panic(err)
	}

	print("Server running on port " + config.GetServerPort())
}
