package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var cotacao Cotacao

	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(cotacao.Usdbrl.Bid)
	if err != nil {
		panic(err)
	}

	err = salvarCotacao(cotacao)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusRequestTimeout)
	}
}

func salvarCotacao(cotacao Cotacao) error {
	dbCtx := context.Background()
	dbCtx, cancel := context.WithTimeout(dbCtx, 10*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite3", "./cotacao.db")

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE Cotacao (
		Id INTEGER PRIMARY KEY AUTOINCREMENT
		,Code       TEXT
		,Codein     TEXT
		,Name       TEXT
		,High       TEXT
		,Low        TEXT
		,VarBid     TEXT
		,PctChange  TEXT
		,Bid        TEXT
		,Ask        TEXT
		,Timestamp  TEXT
		,CreateDate TEXT
		)`)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`INSERT INTO Cotacao(
		Code
		,Codein
		,Name
		,High
		,Low
		,VarBid
		,PctChange
		,Bid
		,Ask
		,Timestamp
		,CreateDate)
		VALUES(?,?,?,?,?,?,?,?,?,?,DATETIME(?))`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(dbCtx,
		cotacao.Usdbrl.Code,
		cotacao.Usdbrl.Codein,
		cotacao.Usdbrl.Name,
		cotacao.Usdbrl.High,
		cotacao.Usdbrl.Low,
		cotacao.Usdbrl.VarBid,
		cotacao.Usdbrl.PctChange,
		cotacao.Usdbrl.Bid,
		cotacao.Usdbrl.Ask,
		cotacao.Usdbrl.Timestamp,
		cotacao.Usdbrl.CreateDate)
	if err != nil {
		return err
	}

	return nil
}
