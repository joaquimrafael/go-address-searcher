package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	IBGE        string `json:"ibge"`
	Gia         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func (a Address) Print() {
	fmt.Printf("Parameter: Cep, Value: %s\n", a.Cep)
	fmt.Printf("Parameter: Logradouro, Value: %s\n", a.Logradouro)
	fmt.Printf("Parameter: Complemento, Value: %s\n", a.Complemento)
	fmt.Printf("Parameter: Unidade, Value: %s\n", a.Unidade)
	fmt.Printf("Parameter: Bairro, Value: %s\n", a.Bairro)
	fmt.Printf("Parameter: Localidade, Value: %s\n", a.Localidade)
	fmt.Printf("Parameter: Uf, Value: %s\n", a.Uf)
	fmt.Printf("Parameter: Estado, Value: %s\n", a.Estado)
	fmt.Printf("Parameter: Regiao, Value: %s\n", a.Regiao)
	fmt.Printf("Parameter: IBGE, Value: %s\n", a.IBGE)
	fmt.Printf("Parameter: Gia, Value: %s\n", a.Gia)
	fmt.Printf("Parameter: DDD, Value: %s\n", a.DDD)
	fmt.Printf("Parameter: Siafi, Value: %s\n", a.Siafi)
}

type result struct {
	addr Address
	err  error
}

func cepClient(ctx context.Context, url string, client http.Client) (Address, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Address{}, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return Address{}, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		return Address{}, fmt.Errorf("invalid url %v", url)
	}

	defer resp.Body.Close()

	var address Address

	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return Address{}, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return address, nil
}

func BuscaCEP(ctx context.Context, cep string, client http.Client) (Address, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	address, err := cepClient(ctx, url, client)
	if err != nil {
		return Address{}, fmt.Errorf("failed to search cep %s: %w", cep, err)
	}

	return address, nil
}

func BuscaVarios(ctx context.Context, client http.Client, ceps []string) ([]Address, error) {
	results := make(chan result)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, v := range ceps {
		wg.Add(1)
		go func() {
			defer wg.Done()
			address, err := BuscaCEP(ctx, v, client)

			select {
			case results <- result{address, err}:
			case <-ctx.Done():
				return
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	out := []Address{}
	for a := range results {
		if a.err != nil {
			cancel()
			return out, a.err
		}
		out = append(out, a.addr)
	}

	return out, nil
}

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ctx := context.Background()

	addresses, err := BuscaVarios(ctx, *client, []string{"08780170", "08710190", "01310100", "20040030", "40026010", "80010000", "69005070"})
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	for _, v := range addresses {
		v.Print()
		fmt.Println()
	}
}
