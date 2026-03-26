package go_cielo_conecta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type (
	Client struct {
		sync.Mutex

		Client   *http.Client
		merchant Merchant
		env      Environment
		token    *tokenResponse
		Log      io.Writer

		cancel context.CancelFunc
		wg     sync.WaitGroup
		once   sync.Once
	}

	tokenResponse struct {
		AccessToken string        `json:"access_token"`
		TokenType   string        `json:"token_type"`
		ExpiresIn   time.Duration `json:"expires_in"`
	}

	Merchant struct {
		ID, Secret string
	}

	ErrorResponse struct {
		Response *http.Response `json:"-"`
		Code     int            `json:",omitempty"`
		Message  string         `json:",omitempty"`
	}

	Sale struct {
		MerchantOrderId string    `json:",omitempty"`
		Customer        *Customer `json:",omitempty"`
		Payment         *Payment  `json:",omitempty"`
	}

	Customer struct {
		Name            string       `json:",omitempty"`
		Identity        string       `json:",omitempty"`
		IdentityType    IdentityType `json:",omitempty"`
		Email           string       `json:",omitempty"`
		Birthday        string       `json:",omitempty"`
		Address         *Address     `json:",omitempty"`
		DeliveryAddress *Address     `json:",omitempty"`
	}

	Address struct {
		Street, Number, Complement, City, State, ZipCode, Country string `json:",omitempty"`
	}

	Payment struct {
		Installments              int                `json:",omitempty"`
		Type                      string             `json:",omitempty"`
		Interest                  Interest           `json:",omitempty"`
		Capture                   bool               `json:",omitempty"`
		SoftDescriptor            string             `json:",omitempty"`
		CreditCard                *CreditCard        `json:",omitempty"`
		PaymentDateTime           string             `json:",omitempty"`
		Amount                    uint64             `json:",omitempty"`
		ProductId                 uint               `json:",omitempty"`
		ReceivedDate              string             `json:",omitempty"`
		CapturedAmount            uint64             `json:",omitempty"`
		CapturedDate              string             `json:",omitempty"`
		Provider                  string             `json:",omitempty"`
		Status                    uint               `json:",omitempty"`
		PhysicalTransactionStatus uint               `json:",omitempty"`
		IsSplitted                bool               `json:",omitempty"`
		ReturnMessage             string             `json:",omitempty"`
		ExtendedMessage           string             `json:",omitempty"`
		ReturnCode                string             `json:",omitempty"`
		PaymentId                 string             `json:",omitempty"`
		Currency                  string             `json:",omitempty"`
		Country                   string             `json:",omitempty"`
		Links                     []Link             `json:",omitempty"`
		ServiceTaxAmount          uint64             `json:",omitempty"`
		PinPadInformation         *PinPadInformation `json:",omitempty"`
		PrintMessage              interface{}        `json:",omitempty"`
		ReceiptInformation        map[string]any     `json:",omitempty"`
		Receipt                   map[string]any     `json:",omitempty"`
		AuthorizationCode         string             `json:",omitempty"`
		ProofOfSale               string             `json:",omitempty"`
		InitializationVersion     string             `json:",omitempty"`
		ConfirmationStatus        uint               `json:",omitempty"`
		EmvResponseData           string             `json:",omitempty"`
		SubordinatedMerchantId    string             `json:",omitempty"`
		OfflinePaymentType        string             `json:",omitempty"`
		MerchantAcquirerId        string             `json:",omitempty"`
		TerminalAcquirerId        string             `json:",omitempty"`
	}

	Link struct {
		Method string `json:",omitempty"`
		Rel    string `json:",omitempty"`
		Href   string `json:",omitempty"`
	}

	CreditCard struct {
		InputMode                      InputMode            `json:",omitempty"`
		ExpirationDate                 string               `json:",omitempty"`
		TrackOneData                   string               `json:",omitempty"`
		TrackTwoData                   string               `json:",omitempty"`
		EncryptedCardData              EncryptedCardData    `json:",omitempty"`
		EmvData                        string               `json:",omitempty"`
		IssuerId                       int                  `json:",omitempty"`
		SecurityCodeStatus             SecurityCodeStatus   `json:",omitempty"`
		SecurityCode                   string               `json:",omitempty"`
		TruncateCardNumberWhenPrinting bool                 `json:",omitempty"`
		SaveCard                       bool                 `json:",omitempty"`
		PanSequenceNumber              uint                 `json:",omitempty"`
		IsFallback                     bool                 `json:",omitempty"`
		AuthenticationMethod           AuthenticationMethod `json:",omitempty"`
		BrandId                        int                  `json:",omitempty"`
		BrandInformation               BrandInformation     `json:",omitempty"`
		PinBlock                       PinBlock             `json:",omitempty"`
	}

	BrandInformation struct {
		Type string `json:",omitempty"`
		Name string `json:",omitempty"`
	}

	PinBlock struct {
		EncryptedPinBlock string         `json:",omitempty"`
		EncryptionType    EncryptionType `json:",omitempty"`
		KsnIdentification string         `json:",omitempty"`
	}

	EncryptedCardData struct {
		EncryptionType       EncryptionType `json:",omitempty"`
		TrackOneDataKSN      string         `json:",omitempty"`
		TrackTwoDataKSN      string         `json:",omitempty"`
		InitializationVector string         `json:",omitempty"`

		IsDataInTLVFormat bool `json:",omitempty"`
	}

	PinPadInformation struct {
		PhysicalCharacteristics PhysicalCharacteristics `json:",omitempty"`
		ReturnDataInfo          string                  `json:",omitempty"`
		SerialNumber            string                  `json:",omitempty"`
		TerminalID              string                  `json:",omitempty"`
	}
)

type (
	IdentityType            string
	Interest                string
	InputMode               string
	AuthenticationMethod    string
	SecurityCodeStatus      string
	EncryptionType          string
	PhysicalCharacteristics string

	currency string
)

type Environment struct {
	OAuthURL     string
	ParamsURL    string
	APIUrl       string
	APIQueryUrl  string
	Homologation bool
}

// Error method implementation for ErrorResponse
func (er ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %3d %s", er.Response.Request.Method, er.Response.Request.URL, er.Response.StatusCode, er.Message)
}
