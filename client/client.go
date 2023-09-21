package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	f, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)

	var cotacao string

	err = json.Unmarshal(body, &cotacao)

	_, err = f.Write([]byte(fmt.Sprintf("Dólar: %s", cotacao)))

	f.Close()

	if err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
		log.Println("Tempo da requisição excedido")
		fmt.Println("Tempo da requisição excedido")
	}
}
