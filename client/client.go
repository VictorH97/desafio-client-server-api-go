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
		log.Println(err)
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	f, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)

	var cotacao string

	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(fmt.Sprintf("Dólar: %s", cotacao)))
	if err != nil {
		panic(err)
	}

	f.Close()
}
