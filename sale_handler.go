package go_cielo_conecta

import (
	"errors"
	"fmt"
)

type SaleInterface interface {
	Authorize() (Sale, error)
	ConfirmPayment(authorizedSale Sale, issuerScriptResults ...string) (ConfirmResponse, error)
	Reverse(confirmedSale Sale, issuerScriptResults ...string) (ConfirmResponse, error)

	SetInstallments(installments int) SaleInterface
	SetInterest(interestType Interest) SaleInterface
	SetCustomer(customer Customer) SaleInterface
	SetPinPadInfo(pinPad PinPadInformation) SaleInterface
	SetSoftDescriptor(softDesc string) SaleInterface

	Get() Sale

	WithCreditCard(cc CreditCard) SaleInterface
	WithDebitCard(dc DebitCard) SaleInterface
}

type SaleHandler struct {
	client *Client
	Sale   Sale
}

func (h *SaleHandler) Get() Sale {
	return h.Sale
}

func newSaleHandler(c *Client, s Sale) SaleInterface {
	return &SaleHandler{client: c, Sale: s}
}

func (h *SaleHandler) WithCreditCard(cc CreditCard) SaleInterface {
	h.Sale.Payment.DebitCard = nil
	h.Sale.Payment.CreditCard = &cc
	h.Sale.Payment.Type = "PhysicalCreditCard"
	return h
}

func (h *SaleHandler) WithDebitCard(dc DebitCard) SaleInterface {
	h.Sale.Payment.CreditCard = nil
	h.Sale.Payment.DebitCard = &dc
	h.Sale.Payment.Type = "PhysicalDebitCard"
	return h
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
// Returns the authorized sale with payment details or an error if the validation fails or if there is an issue with the API requestBody.
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

	if salePayed.Payment != nil && salePayed.Payment.Status != StatusPaymentConfirmed {
		return salePayed, ErrPaymentIsNotConfirmed
	}

	return salePayed, nil
}

// ConfirmPayment confirms a payment with the provided issuer script results.
// Returns the confirmation result or an error if the validation fails or if there is an issue with the API requestBody.
//
// PUT /1/physicalSales/{PaymentId}/confirmation
func (h *SaleHandler) ConfirmPayment(authorizedSale Sale, issuerScriptResults ...string) (ConfirmResponse, error) {
	if authorizedSale.Payment == nil {
		return ConfirmResponse{}, ErrPaymentRequired
	}

	var response ConfirmResponse

	href := fmt.Sprintf("%s/1/physicalSales/%s/confirmation", h.client.env.APIUrl, authorizedSale.Payment.PaymentId)

	body := map[string]string{
		"EmvData":             h.Sale.Payment.getEmvData(),
		"IssuerScriptResults": "0000",
	}

	if len(issuerScriptResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptResults[0]
	}

	req, err := h.client.NewRequest("PUT", href, body)
	if err != nil {
		return response, err
	}

	err = h.client.Send(req, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// ReversePayment must be called when payment returns success. It initiates the reversal process for a given sale, allowing for the cancellation of a previously authorized payment.
// The method accepts a Sale object and an optional issuerScriptsResults string, which can be used to provide additional information for the reversal process.
//
// Depending on whether the Sale object contains a PaymentId, the method will choose the appropriate reversal endpoint (either by PaymentId or by OrderId) and send a DELETE requestBody to the API.
// It returns a ConfirmResponse indicating the result of the reversal operation or an error if the requestBody fails.
func (h *SaleHandler) Reverse(authorizedSale Sale, issuerScriptResults ...string) (ConfirmResponse, error) {
	emvData := h.Sale.Payment.getEmvData()

	cancel := newCancelHandler(h.client, emvData, issuerScriptResults...)

	h.client.LogInfo("trying to reverse payment", "payment", authorizedSale.Payment)

	return cancel.ReverseWithPaymentID(authorizedSale.Payment.PaymentId)
}

func (h *SaleHandler) validate() error {
	var errs error

	if h.Sale.MerchantOrderId == "" {
		errs = errors.Join(errs, ErrOrderIDRequired)
	}

	if h.Sale.Payment.Type == "" {
		errs = errors.Join(errs, ErrPaymentTypeRequired)
	}

	if h.Sale.Payment.SoftDescriptor == "" {
		errs = errors.Join(errs, ErrSoftDescriptorRequired)
	}

	if h.Sale.Payment.CreditCard == nil && h.Sale.Payment.DebitCard == nil {
		errs = errors.Join(errs, ErrCardRequired)
	}

	return errs
}
