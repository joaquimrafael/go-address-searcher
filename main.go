package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func cepClient(url string) (Address, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

func BuscaCEP(cep string) (Address, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	address, err := cepClient(url)
	if err != nil {
		return Address{}, fmt.Errorf("failed to search cep %s: %w", cep, err)
	}

	return address, nil
}

func BuscaVarios(ctx context.Context, ceps []string) ([]Address, error) {
	addresses := make([]Address, 0, len(ceps))
	for _, v := range ceps {
		address, err := BuscaCEP(v)
		if err != nil {
			return addresses, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func main() {
	address, err := BuscaCEP("01001000")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	fmt.Println("Address:", address)

	addresses, err := BuscaVarios(context.TODO(), []string{"08780170", "08710190"})
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	fmt.Println("Addresses:", addresses)
}
