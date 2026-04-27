package go_cielo_conecta

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type (
	Client struct {
		sync.Mutex

		Client *http.Client
		env    Environment
		token  *tokenResponse
		log    *slog.Logger

		cancel context.CancelFunc
		wg     sync.WaitGroup
		once   sync.Once
	}

	Environment struct {
		OAuthURL     string
		ParamsURL    string
		APIUrl       string
		APIQueryUrl  string
		Homologation bool
		merchant     Merchant
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

	MultiError struct {
		Errors []ErrorResponse `json:"errors"`
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
		Installments              int                   `json:",omitempty"`
		Type                      string                `json:",omitempty"`
		Interest                  Interest              `json:",omitempty"`
		Capture                   bool                  `json:",omitempty"` // Capture identifica que a autorização deve ser com captura automática. A autorização sem captura automática é conhecida também como pré-autorização.
		SoftDescriptor            string                `json:",omitempty"`
		CreditCard                *CreditCard           `json:",omitempty"`
		DebitCard                 *DebitCard            `json:",omitempty"`
		PaymentDateTime           string                `json:",omitempty"`
		Amount                    uint64                `json:",omitempty"`
		ProductId                 uint                  `json:",omitempty"`
		ReceivedDate              string                `json:",omitempty"`
		CapturedAmount            uint64                `json:",omitempty"`
		CapturedDate              string                `json:",omitempty"`
		Provider                  string                `json:",omitempty"`
		Status                    StatusPayment         `json:",omitempty"`
		PhysicalTransactionStatus uint                  `json:",omitempty"`
		IsSplitted                bool                  `json:",omitempty"`
		ReturnMessage             string                `json:",omitempty"`
		ExtendedMessage           string                `json:",omitempty"`
		ReturnCode                string                `json:",omitempty"`
		PaymentId                 string                `json:",omitempty"`
		Currency                  string                `json:",omitempty"`
		Country                   string                `json:",omitempty"`
		Links                     []*Link               `json:",omitempty"`
		ServiceTaxAmount          uint64                `json:",omitempty"`
		PinPadInformation         *PinPadInformation    `json:",omitempty"`
		PrintMessage              interface{}           `json:",omitempty"`
		ReceiptInformation        []*ReceiptInformation `json:",omitempty"`
		Receipt                   map[string]any        `json:",omitempty"`
		AuthorizationCode         string                `json:",omitempty"`
		ProofOfSale               string                `json:",omitempty"`
		InitializationVersion     int64                 `json:",omitempty"`
		ConfirmationStatus        uint                  `json:",omitempty"`
		EmvResponseData           string                `json:",omitempty"`
		SubordinatedMerchantId    string                `json:",omitempty"`
		OfflinePaymentType        string                `json:",omitempty"`
		MerchantAcquirerId        string                `json:",omitempty"`
		TerminalAcquirerId        string                `json:",omitempty"`
	}

	ReceiptInformation struct {
		Field   string `json:",omitempty"`
		Label   string `json:",omitempty"`
		Content string `json:",omitempty"`
	}

	Link struct {
		Method string `json:",omitempty"`
		Rel    string `json:",omitempty"`
		Href   string `json:",omitempty"`
	}

	CreditCard struct {
		InputMode                      InputMode            `json:",omitempty"`
		ExpirationDate                 string               `json:",omitempty"`
		AuthenticationMethod           AuthenticationMethod `json:",omitempty"`
		IssuerId                       int                  `json:",omitempty"`
		BrandId                        int                  `json:",omitempty"`
		TrackOneData                   string               `json:",omitempty"`
		TrackTwoData                   string               `json:",omitempty"`
		EmvData                        string               `json:",omitempty"`
		EncryptedCardData              EncryptedCardData    `json:",omitempty"`
		SecurityCodeStatus             SecurityCodeStatus   `json:",omitempty"`
		SecurityCode                   string               `json:",omitempty"`
		TruncateCardNumberWhenPrinting bool                 `json:",omitempty"`
		SaveCard                       bool                 `json:",omitempty"`
		PanSequenceNumber              uint                 `json:",omitempty"`
		IsFallback                     bool                 `json:",omitempty"`
		BrandInformation               BrandInformation     `json:",omitempty"`
		PinBlock                       PinBlock             `json:",omitempty"`
	}

	DebitCard struct {
		InputMode                      InputMode            `json:",omitempty"`
		ExpirationDate                 string               `json:",omitempty"`
		AuthenticationMethod           AuthenticationMethod `json:",omitempty"`
		IssuerId                       uint                 `json:",omitempty"`
		BrandId                        uint                 `json:",omitempty"`
		TruncateCardNumberWhenPrinting bool                 `json:",omitempty"`
		PanSequenceNumber              uint                 `json:",omitempty"`
		SaveCard                       bool                 `json:",omitempty"`
		EmvData                        string               `json:",omitempty"`
		TrackOneData                   string               `json:",omitempty"`
		TrackTwoData                   string               `json:",omitempty"`
		EncryptedCardData              EncryptedCardData    `json:",omitempty"`
		PinBlock                       PinBlock             `json:",omitempty"`
		IsFallback                     bool                 `json:",omitempty"`
		CardToken                      string               `json:",omitempty"`
		BrandInformation               BrandInformation     `json:",omitempty"`
		SecurityCodeStatus             SecurityCodeStatus   `json:",omitempty"`
		SecurityCode                   string               `json:",omitempty"`
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
		IsDataInTLVFormat    bool           `json:",omitempty"`
	}

	PinPadInformation struct {
		PhysicalCharacteristics PhysicalCharacteristics `json:",omitempty"`
		ReturnDataInfo          string                  `json:",omitempty"`
		SerialNumber            string                  `json:",omitempty"`
		TerminalID              string                  `json:",omitempty"`
	}

	ConfirmResponse struct {
		ConfirmationStatus ConfirmationStatus `json:",omitempty"`
		Status             uint16             `json:",omitempty"`
		ReasonCode         uint16             `json:",omitempty"`
		ReturnCode         string             `json:",omitempty"`
		ReturnMessage      string             `json:",omitempty"`
		Links              []*Link            `json:",omitempty"`
	}
)
type (
	IdentityType            string
	Interest                string
	InputMode               string
	AuthenticationMethod    string
	SecurityCodeStatus      string
	PhysicalCharacteristics string

	currency string
)

type (
	EncryptionType     uint
	StatusPayment      uint
	ConfirmationStatus uint
)

func (e Environment) WithMerchant(m Merchant) Environment {
	e.merchant = m
	return e
}

// Error method implementation for ErrorResponse
func (er ErrorResponse) Error() string {
	return fmt.Sprintf("%s %s: %3d %s", er.Response.Request.Method, er.Response.Request.URL, er.Response.StatusCode, er.Message)
}

func (er MultiError) Error() string {
	var msgs []string
	for _, err := range er.Errors {
		msgs = append(msgs, err.Error())
	}

	return strings.Join(msgs, "; ")
}

func (p Payment) getLink(rel string) *Link {
	if len(p.Links) == 0 {
		return nil
	}

	for _, l := range p.Links {
		if l.Rel == rel {
			return l
		}
	}

	return nil
}

func (p Payment) getEmvData() string {
	if p.CreditCard != nil {
		return p.CreditCard.EmvData
	}

	if p.DebitCard != nil {
		return p.CreditCard.EmvData
	}

	return ""
}
