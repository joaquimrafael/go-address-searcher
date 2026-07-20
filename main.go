package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	cep := "01001000"
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	log.Printf("Status: %s", resp.Status)
	log.Printf("Body: %s", string(body))
}
