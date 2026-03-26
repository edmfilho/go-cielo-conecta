# Go Cielo Conecta
___
Golang lib to communicate with Cielo Conecta API.

### Prerequisite Knowledge
- <a href="https://developercielo.github.io/manual/cielo-conecta" target="_blank">Docs Cielo Conecta</a> 

### Requirements
- Go 1.18 or higher
- Cielo Conecta API credentials

### TODOs:
- Better tests
- Implement requests logs
- Implement endpoints:
  - Terminals
  - PaymentsQuery
  - Cancellation
  - Stores
  - Equipments

### Installation

```shell
go get github.com/edmfilho/go-cielo-conecta
```

### New Client

```go
package main

import "github.com/edmfilho/cieloconecta"

func main() {
	merchant := cieloconecta.Merchant{
		ID:     "your_merchant_id",
		Secret: "your_merchant_secret",
	}

	cieloClient, err := cieloconecta.NewClient(merchant, cieloconecta.SandboxEnvironment)
	if err != nil {
		panic(err)
	}
  
	// Use cieloClient to make API calls here
	defer cieloClient.Close()
}

```



