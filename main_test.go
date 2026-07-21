package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuscaCEP(t *testing.T) {
	cep := "80010000"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"80010-000","localidade":"Curitiba","uf":"PR"}`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	originalUrl := urlViaCep
	urlViaCep = server.URL + "/"
	defer func() { urlViaCep = originalUrl }()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	ctx := context.Background()
	address, err := BuscaCEP(ctx, cep, *client)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if address == (Address{}) {
		t.Errorf("expected a filled address, got empty")
	}

	if address.Localidade != "Curitiba" {
		t.Errorf("expected Localidade %q, got %q", "Curitiba", address.Localidade)
	}
}

func TestBuscaVarios(t *testing.T) {
	ceps := []string{"80010000", "80010000"}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"80010-000","localidade":"Curitiba","uf":"PR"}`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	originalUrl := urlViaCep
	urlViaCep = server.URL + "/"
	defer func() { urlViaCep = originalUrl }()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	ctx := context.Background()
	addresses, err := BuscaVarios(ctx, *client, ceps)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(addresses) != 2 {
		t.Errorf("expected a filled slice, got empty")
	}
}
