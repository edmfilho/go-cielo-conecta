package go_cielo_conecta

import (
	"errors"
	"fmt"
	"log/slog"
)

type SaleInterface interface {
	Authorize() (Sale, error)
	ConfirmPayment(authorizedSale Sale, issuerScriptResults ...string) (ConfirmResponse, error)

	SetInstallments(installments int) SaleInterface
	SetInterest(interestType Interest) SaleInterface
	SetCustomer(customer Customer) SaleInterface
	SetPinPadInfo(pinPad PinPadInformation) SaleInterface
	SetSoftDescriptor(softDesc string) SaleInterface

	WithCreditCard(cc CreditCard) SaleInterface
	WithDebitCard(dc DebitCard) SaleInterface
}

type SaleHandler struct {
	client *Client
	Sale   Sale
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

	if authorizedSale.Payment.Status != StatusPaymentConfirmed {
		h.client.LogError(fmt.Sprintf("[%s] payment status is not confirmed", authorizedSale.Payment.PaymentId),
			slog.Group("payment",
				slog.String("status", authorizedSale.Payment.Status.String()),
				slog.String("message", authorizedSale.Payment.ExtendedMessage),
				slog.String("return_message", authorizedSale.Payment.ReturnMessage),
			),
		)
		return ConfirmResponse{}, fmt.Errorf("payment status is not confirmed: status=%s", authorizedSale.Payment.Status)
	}

	link := authorizedSale.Payment.getLink("confirm")
	if link == nil {
		h.client.LogError(fmt.Sprintf("[%s] could not confirm this payment", authorizedSale.Payment.PaymentId),
			slog.Group("payment",
				slog.String("status", authorizedSale.Payment.Status.String()),
				slog.String("message", authorizedSale.Payment.ExtendedMessage),
				slog.String("return_message", authorizedSale.Payment.ReturnMessage),
			),
		)

		return ConfirmResponse{}, fmt.Errorf("could not confirm this payment, status=%s", authorizedSale.Payment.Status)
	}

	var body = map[string]string{}

	body["EmvData"] = h.Sale.Payment.getEmvData()
	body["IssuerScriptResults"] = "0000"
	if len(issuerScriptResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptResults[0]
	}

	req, err := h.client.NewRequest(link.Method, link.Href, body)
	if err != nil {
		return ConfirmResponse{}, err
	}

	var resp ConfirmResponse

	err = h.client.Send(req, &resp)
	if err != nil {
		return ConfirmResponse{}, err
	}

	return resp, nil

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
