# go-address-searcher

Looks up addresses from Brazilian postal codes (CEP) using the public [ViaCEP](https://viacep.com.br/) API, querying multiple CEPs concurrently.

## Requirements

- Go 1.26 or newer

## Installation

```sh
git clone https://github.com/joaquimrafael/go-address-searcher.git
cd go-address-searcher
```

## Usage

The list of CEPs to look up is defined in `main.go`. Run it with:

```sh
go run .
```

To build the binary:

```sh
go build -o go-address-searcher .
./go-address-searcher
```

Each address found is printed to the terminal:

```json
Cep: 80010-000
Logradouro: Praça Generoso Marques
Complemento:
...
Localidade: Curitiba
Uf: PR
```

## API

The exported functions can be reused as a package:

### `BuscaCEP`

```go
func BuscaCEP(ctx context.Context, cep string, client http.Client) (Address, error)
```

Looks up a single CEP and returns the matching `Address`.

### `BuscaVarios`

```go
func BuscaVarios(ctx context.Context, client http.Client, ceps []string) ([]Address, error)
```

Looks up multiple CEPs concurrently. If any lookup fails, it cancels the in-flight requests and returns the error.

## Tests

```sh
go test ./...
```

With the race detector:

```sh
go test -race ./...
```

The tests use `net/http/httptest` to simulate the ViaCEP API, without making real network calls.
