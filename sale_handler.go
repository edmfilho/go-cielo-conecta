package go_cielo_conecta

import (
	"errors"
	"fmt"
)

type SaleInterface interface {
	GetSale() *Sale
	Exec() (*Sale, error)

	WithCreditCardOnlinePassword(cc *CreditCard) SaleInterface

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
	return &SaleHandler{client: client, Sale: s}
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

func (h *SaleHandler) GetSale() *Sale {
	return h.Sale
}

func (h *SaleHandler) SetInstallments(installments int) SaleInterface {
	h.Sale.Payment.Installments = installments
	return h
}

func (h *SaleHandler) WithCreditCardOnlinePassword(cc *CreditCard) SaleInterface {
	h.Sale.Payment.CreditCard = cc
	h.Sale.Payment.Type = "PhysicalCreditCard"
	return h
}

// POST /1/physicalSales/
func (h *SaleHandler) Exec() (*Sale, error) {
	salePayed := &Sale{}

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
