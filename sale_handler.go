package go_cielo_conecta

import (
	"errors"
	"fmt"
	"slices"
)

type SaleInterface interface {
	GetSale() Sale
	Authorization() (Sale, error)
	Confirm(emvData string, issuerScriptResults ...string) (*ConfirmResponse, error)

	WithCreditCard(cc *CreditCard) SaleInterface
	WithDebitCard(dc *DebitCard) SaleInterface

	SetInstallments(installments int) SaleInterface
	SetInterest(interestType Interest) SaleInterface
	SetCustomer(customer *Customer) SaleInterface
	SetPinPadInfo(pinPad *PinPadInformation) SaleInterface
	SetSoftDescriptor(softDesc string) SaleInterface
}

type SaleHandler struct {
	client *Client
	Sale   *Sale
}

func newSaleHandler(client *Client, s *Sale) SaleInterface {
	return &SaleHandler{
		client: client,
		Sale: &Sale{
			MerchantOrderId: s.MerchantOrderId,
			Payment:         s.Payment,
		},
	}
}

func (h *SaleHandler) SetSoftDescriptor(softDesc string) SaleInterface {
	h.Sale.Payment.SoftDescriptor = softDesc
	return h
}

func (h *SaleHandler) SetPinPadInfo(pinPad *PinPadInformation) SaleInterface {
	h.Sale.Payment.PinPadInformation = pinPad
	return h
}

func (h *SaleHandler) SetCustomer(customer *Customer) SaleInterface {
	h.Sale.Customer = customer
	return h
}

func (h *SaleHandler) SetInterest(interestType Interest) SaleInterface {
	h.Sale.Payment.Interest = interestType
	return h
}

func (h *SaleHandler) GetSale() Sale {
	return *h.Sale
}

func (h *SaleHandler) SetInstallments(installments int) SaleInterface {
	h.Sale.Payment.Installments = installments
	return h
}

func (h *SaleHandler) WithCreditCard(cc *CreditCard) SaleInterface {
	h.Sale.Payment.CreditCard = cc
	h.Sale.Payment.DebitCard = nil
	h.Sale.Payment.Type = "PhysicalCreditCard"
	return h
}

func (h *SaleHandler) WithDebitCard(dc *DebitCard) SaleInterface {
	h.Sale.Payment.CreditCard = nil
	h.Sale.Payment.DebitCard = dc
	h.Sale.Payment.Type = "PhysicalDebitCard"
	return h
}

// Authorization validates the sale data and sends a request to the API to authorize the payment.
// It returns the authorized sale or an error if the validation fails or if there is an issue with the API request.
// POST /1/physicalSales/
func (h *SaleHandler) Authorization() (Sale, error) {
	salePayed := Sale{}

	if err := h.validate(); err != nil {
		return salePayed, err
	}

	req, err := h.client.NewRequest("POST", fmt.Sprintf("%s%s", h.client.env.APIUrl, "/1/physicalSales/"), h.Sale)
	if err != nil {
		return salePayed, err
	}

	err = h.client.Send(req, &salePayed)
	if err != nil {
		return salePayed, err
	}

	h.Sale = &salePayed

	return salePayed, nil
}

// Confirm confirms a payment with the provided issuer script results.
// It validates the sale data and sends a request to the API to confirm the payment.
// It returns the confirmation result or an error if the validation fails or if there is an issue with the API request.
//
// PUT /1/physicalSales/{PaymentId}/confirmation
func (h *SaleHandler) Confirm(emvData string, issuerScriptResults ...string) (result *ConfirmResponse, err error) {
	if h.Sale == nil {
		return nil, errors.New("sale not initialized")
	}

	if h.Sale.Payment.PaymentId == "" {
		return nil, errors.New("payment_id is required")
	}

	if h.Sale.Payment.Status != Confirmed {
		return nil, fmt.Errorf("payment is not confirmed: status=%s", h.Sale.Payment.Status)
	}

	i := slices.IndexFunc(h.Sale.Payment.Links, func(link *Link) bool { return link.Rel == "confirm" })

	var (
		urlConfirm = h.Sale.Payment.Links[i].Href
		method     = h.Sale.Payment.Links[i].Method
		body       = make(map[string]string)
	)

	body["EmvData"] = emvData
	body["IssuerScriptResults"] = "0000"

	if len(issuerScriptResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptResults[0]
	}

	req, err := h.client.NewRequest(method, urlConfirm, body)
	if err != nil {
		return result, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (h *SaleHandler) validate() error {
	var errs error

	if h.Sale.MerchantOrderId == "" {
		errs = errors.New("merchant order id is required")
	}

	if h.Sale.Payment.Type == "" {
		errs = errors.Join(errs, errors.New("payment type is required"))
	}

	if h.Sale.Payment.SoftDescriptor == "" {
		errs = errors.Join(errs, errors.New("soft descriptor is required"))
	}

	if h.Sale.Payment.CreditCard == nil {
		errs = errors.Join(errs, errors.New("no credit card"))
	}

	return errs
}
