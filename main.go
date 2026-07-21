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

func cepClient(url string, client http.Client) (Address, error) {
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

func BuscaCEP(cep string, client http.Client) (Address, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	address, err := cepClient(url, client)
	if err != nil {
		return Address{}, fmt.Errorf("failed to search cep %s: %w", cep, err)
	}

	return address, nil
}

func BuscaVarios(ctx context.Context, client http.Client, ceps []string) ([]Address, error) {
	addresses := make([]Address, 0, len(ceps))
	for _, v := range ceps {
		address, err := BuscaCEP(v, client)
		if err != nil {
			return addresses, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	address, err := BuscaCEP("01001000", *client)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	address.Print()

	addresses, err := BuscaVarios(context.TODO(), *client, []string{"08780170", "08710190"})
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	fmt.Println()

	for _, v := range addresses {
		v.Print()
		fmt.Println()
	}
}
