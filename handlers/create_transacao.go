package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thiagosousasantana/rinha-go/db"
	"github.com/thiagosousasantana/rinha-go/models"
	"github.com/thiagosousasantana/rinha-go/repositories"
)

func CreateTransacao(w http.ResponseWriter, r *http.Request) {
	var transacao models.Transacao

	clienteId, err := strconv.Atoi(chi.URLParam(r, "id"))

	err = json.NewDecoder(r.Body).Decode(&transacao)
	if err != nil {
		log.Printf("Error docoding request response. %v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = ValidateTransacao(transacao, clienteId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	transacao.Data = time.Now().Format("2006-01-02T15:04:05.999999999Z07:00")
	transacao.ClienteID = clienteId

	ctx := r.Context()
	tx, err := db.OpenTransaction(ctx)
	if err != nil {
		log.Printf("Error opening transaction. %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cliente, err := repositories.FindByIdLocking(clienteId, tx)
	if err != nil {
		log.Printf("Error getting client. %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Printf("Rollback error: %v", rbErr)
		}
		return
	}

	if transacao.Tipo == "d" {
		novoSaldo := cliente.Saldo - transacao.Valor
		if novoSaldo+cliente.Limite < 0 {
			http.Error(w, "Insufficient funds", http.StatusUnprocessableEntity)
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Printf("Rollback error: %v", rbErr)
			}
			return
		}
		cliente.Saldo = novoSaldo
	} else {
		cliente.Saldo += transacao.Valor
	}

	_, err = repositories.UpdateSaldo(int8(cliente.ID), cliente.Saldo, tx)
	if err != nil {
		log.Printf("Error updating client balance. %v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = repositories.InsertTransacao(transacao)
	if err != nil {
		log.Printf("Error creating new transacao. %v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	var resp map[string]any

	if err != nil {
		resp = map[string]any{
			"status":  "error",
			"message": fmt.Sprintf("Error creating new transacao. %v", err),
		}
	} else {
		resp = map[string]any{
			"limite": cliente.Limite,
			"saldo":  cliente.Saldo,
		}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ValidateTransacao(transacao models.Transacao, clienteId int) error {
	if transacao.Tipo != "d" && transacao.Tipo != "c" {
		return fmt.Errorf("Invalid transaction type. Type: %s", transacao.Tipo)
	}

	if transacao.Valor < 1 {
		return fmt.Errorf("Invalid transaction value")
	}

	if transacao.Descricao == "" || len(transacao.Descricao) > 10 {
		return fmt.Errorf("Invalid transaction description")
	}
	return nil
}
