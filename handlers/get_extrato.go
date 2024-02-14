package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thiagosousasantana/rinha-go/repositories"
)

func GetExtrato(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Error converting id to integer. %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ValidateId(id)
	if err != nil {
		log.Printf("Cliente not found id. %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cliente, err := repositories.FindById(id)
	if err != nil {
		log.Printf("Error getting client. %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	saldo := map[string]any{
		"total":        cliente.Saldo,
		"data_extrato": time.Now().Format("2024-01-17T02:34:41.217753Z"),
		"limite":       cliente.Limite,
	}

	transacoes, err := repositories.FindLast10ByClienteId(int8(id))
	if err != nil {
		log.Printf("Error getting transactions. %v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	ultimasTransacoes := []map[string]any{}

	for _, transacao := range transacoes {
		ultimaTransacao := map[string]any{
			"valor":        transacao.Valor,
			"tipo":         transacao.Tipo,
			"descricao":    transacao.Descricao,
			"realizada_em": transacao.Data,
		}
		ultimasTransacoes = append(ultimasTransacoes, ultimaTransacao)
	}

	resp := map[string]any{
		"saldo":              saldo,
		"ultimas_transacoes": ultimasTransacoes,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ValidateId(id int) error {
	if id < 1 || id > 5 {
		return errors.New("Invalid id")
	}
	return nil
}
