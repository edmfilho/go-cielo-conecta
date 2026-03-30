# Go Cielo Conecta
___
Golang lib to communicate with Cielo Conecta API.

### Prerequisite knowledge
- <a href="https://developercielo.github.io/manual/cielo-conecta" target="_blank">Docs Cielo Conecta</a> 

### Requirements
- Go 1.18 or higher
- Cielo Conecta API credentials

### TODOs:
- Better tests
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
var merchant = cieloConecta.Merchant{
  ID:     "your_merchant_id",
  Secret: "your_merchant_secret",
}

// Use client to make API calls
client, err := cieloConecta.NewClient(merchant, cieloConecta.SandboxEnvironment)
if err != nil {
  log.Fatal(err)
}

// Remember to close the client when you're done (this will stop the goroutine that refreshes the token)
defer client.Close()

// By default, the requests will be logged to the standard output. You can disable this by setting the logger to nil.
client.SetLogger(nil)

// By the way, you can also set a custom logger that implements the Logger interface.
client.SetLogger(&MyCustomLogger{})
```

### Create a new payment:

```go
// Read more about creating payments in the documentation: https://developercielo.github.io/manual/cielo-conecta#fluxo-de-pagamento
cc := &cieloConecta.CreditCard{}

sale, err := client.CreatePayment("123456789", 100.0, 1).
    WithCreditCardOnlinePassword(cc).
    SetSoftDescriptor("My Store").
    Exec()
if err != nil {
    log.Fatal(err)
}

log.Println("Payment created successfully: ", sale)
```


