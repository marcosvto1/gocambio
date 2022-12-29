package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CotacaoValue struct {
	Code string `json:"code"`
	CodeIn string `json:"codein"`
	Name string `json:"name"`
	High string `json:"high"`
	Low string `json:"low"`
	VarBid string `json:"varBid"`
	PctChange string `json:"pctChange"`
	Bid string `json:"bid"`
	Ask string `json:"Ask"`
	Timestamp string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacao struct {
	USDBRL CotacaoValue
}

type Server struct {
	DB *sql.DB
}

func main() {
	db, err := sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := Server{DB: db}

	http.HandleFunc("/cotacao", server.ConsultarCotacaoDolarHandle)
	http.ListenAndServe(":8081", nil)
}

func (s *Server) ConsultarCotacaoDolarHandle(w http.ResponseWriter, r *http.Request) {
	cotacao, err := ConsultarCotacao()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	err = SalvarCotacao(s, cotacao)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&cotacao)
}

func SalvarCotacao(s *Server, cotacao *Cotacao) (error) {
	fmt.Println("* Salvando Cotação.")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Millisecond)
	defer cancel()

	stmt, err := s.DB.PrepareContext(ctx,"INSERT INTO cotacoes (result,data_cotacao) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer stmt.Close()

	res, err := json.Marshal(cotacao)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = stmt.ExecContext(ctx, res, time.Now())
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func ConsultarCotacao() (*Cotacao, error) {
	fmt.Println("* Consultando Cotação.")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200 * time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}