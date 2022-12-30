package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CotacaoValue struct {
	Bid string `json:"bid"`
}

type Cotacao struct {
	USDBRL CotacaoValue
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond * 300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	arquivo, err := os.OpenFile("cotacao.txt", os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dólar: %v \n", cotacao.USDBRL.Bid)

	_, err = arquivo.WriteString("Dólar: "+ cotacao.USDBRL.Bid)
	arquivo.Close()
	if err != nil {
		panic(err)
	}
}