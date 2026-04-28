package go_cielo_conecta

import (
	"errors"
	"fmt"
)

type SaleInterface interface {
	Authorization() (Sale, error)

	WithCreditCard(cc *CreditCard) SaleInterface
	WithDebitCard(dc *DebitCard) SaleInterface

	SetInstallments(installments int) SaleInterface
	SetInterest(interestType Interest) SaleInterface
	SetCustomer(customer Customer) SaleInterface
	SetPinPadInfo(pinPad PinPadInformation) SaleInterface
	SetSoftDescriptor(softDesc string) SaleInterface
}

type SaleHandler struct {
	client *Client
	Sale   Sale
}

func newSaleHandler(c *Client, s Sale) SaleInterface {
	return &SaleHandler{client: c, Sale: s}
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

// Authorization validates the sale data and sends a requestBody to the API to authorize the payment.
// It returns the authorized sale or an error if the validation fails or if there is an issue with the API requestBody.
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
