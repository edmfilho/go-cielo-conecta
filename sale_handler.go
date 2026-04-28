package go_cielo_conecta

import (
	"errors"
	"fmt"
)

type SaleInterface interface {
	Authorize() (Sale, error)

	SetInstallments(installments int) SaleInterface
	SetInterest(interestType Interest) SaleInterface
	SetCustomer(customer Customer) SaleInterface
	SetPinPadInfo(pinPad PinPadInformation) SaleInterface
	SetSoftDescriptor(softDesc string) SaleInterface

	WithCard(card any) (SaleInterface, error)
}

type SaleHandler struct {
	client *Client
	Sale   Sale
}

func newSaleHandler(c *Client, s Sale, card any) (SaleInterface, error) {
	h := SaleHandler{client: c, Sale: s}
	return h.WithCard(card)
}

func (h *SaleHandler) WithCard(card any) (SaleInterface, error) {
	if card == nil {
		return h, errors.New("card is required")
	}

	h.Sale.Payment.DebitCard = nil
	h.Sale.Payment.CreditCard = nil

	switch v := card.(type) {
	case CreditCard:
		h.Sale.Payment.CreditCard = &v
		h.Sale.Payment.Type = "PhysicalCreditCard"
	case DebitCard:
		h.Sale.Payment.DebitCard = &v
		h.Sale.Payment.Type = "PhysicalDebitCard"
	default:
		return h, errors.New("card must be of type CreditCard or DebitCard")
	}

	return h, nil
}

func (h *SaleHandler) SetSoftDescriptor(softDesc string) SaleInterface {
	h.Sale.Payment.SoftDescriptor = softDesc
	return h
}

func (h *SaleHandler) SetPinPadInfo(pinPad PinPadInformation) SaleInterface {
	h.Sale.Payment.PinPadInformation = &pinPad
	return h
}

func (h *SaleHandler) SetCustomer(c Customer) SaleInterface {
	h.Sale.Customer = &c
	return h
}

func (h *SaleHandler) SetInterest(interestType Interest) SaleInterface {
	h.Sale.Payment.Interest = interestType
	return h
}

func (h *SaleHandler) SetInstallments(installments int) SaleInterface {
	h.Sale.Payment.Installments = installments
	return h
}

// Authorize validates the sale data and sends a requestBody to the API to authorize the payment.
// It returns the authorized sale or an error if the validation fails or if there is an issue with the API requestBody.
func (h *SaleHandler) Authorize() (Sale, error) {
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

	h.Sale = Sale{}

	return salePayed, nil
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
