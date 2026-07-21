package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setupTest(t *testing.T, handler http.HandlerFunc) *http.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	originalUrl := urlViaCep
	urlViaCep = server.URL + "/"
	t.Cleanup(func() { urlViaCep = originalUrl })

	return &http.Client{
		Timeout: 1 * time.Second,
	}
}

func TestBuscaCEP(t *testing.T) {
	client := setupTest(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"80010-000","localidade":"Curitiba","uf":"PR"}`))
	})

	address, err := BuscaCEP(context.Background(), "80010000", *client)
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
	client := setupTest(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"80010-000","localidade":"Curitiba","uf":"PR"}`))
	})

	ceps := []string{"80010000", "80010000"}

	addresses, err := BuscaVarios(context.Background(), *client, ceps)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(addresses) != len(ceps) {
		t.Errorf("expected %d addresses, got %d", len(ceps), len(addresses))
	}
}

func TestBuscaCEP_StatusNotOK(t *testing.T) {
	client := setupTest(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := BuscaCEP(context.Background(), "80010000", *client)
	if err == nil {
		t.Fatal("expected an error for non-200 status, got nil")
	}
}

func TestBuscaCEP_InvalidJSON(t *testing.T) {
	client := setupTest(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{not a valid json`))
	})

	_, err := BuscaCEP(context.Background(), "80010000", *client)
	if err == nil {
		t.Fatal("expected a decode error for invalid JSON, got nil")
	}
}

func TestBuscaCEP_CanceledContext(t *testing.T) {
	client := setupTest(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"localidade":"Curitiba"}`))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := BuscaCEP(ctx, "80010000", *client)
	if err == nil {
		t.Fatal("expected an error for a canceled context, got nil")
	}
}

func TestBuscaVarios_Error(t *testing.T) {
	client := setupTest(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "99999999") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"80010-000","localidade":"Curitiba","uf":"PR"}`))
	})

	ceps := []string{"80010000", "99999999", "80010000"}

	_, err := BuscaVarios(context.Background(), *client, ceps)
	if err == nil {
		t.Fatal("expected an error when one cep fails, got nil")
	}
}
