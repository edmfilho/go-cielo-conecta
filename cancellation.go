package go_cielo_conecta

import (
	"errors"
	"fmt"
)

type CancelInterface interface {
	ReverseWithPaymentID() (ConfirmResponse, error)
	ReverseWithOrderID() (ConfirmResponse, error)
}

type CancelHandler struct {
	client       *Client
	sale         Sale
	requestBody  map[string]string
	urlPaymentID string
}

func newCancelHandler(c *Client, s Sale, issuerScriptsResults ...string) (CancelInterface, error) {
	if s.Payment == nil {
		return nil, ErrPaymentRequired
	}

	if s.Payment.CreditCard == nil && s.Payment.DebitCard == nil {
		return nil, errors.New("Sale.Payment.CreditCard or Sale.Payment.DebitCard is required")
	}

	var body map[string]string

	link := s.Payment.getLink("reverse")
	if link == nil {
		return nil, fmt.Errorf("could not reverse this payment, status=%s, message=%s %s", s.Payment.Status, s.Payment.ExtendedMessage, s.Payment.ReturnMessage)
	}

	body["EmvData"] = s.Payment.getEmvData()
	body["IssuerScriptResults"] = "0000"

	if len(issuerScriptsResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptsResults[0]
	}

	return &CancelHandler{client: c, sale: s, requestBody: body, urlPaymentID: link.Href}, nil
}

func (h *CancelHandler) ReverseWithPaymentID() (ConfirmResponse, error) {
	var result ConfirmResponse

	req, err := h.client.NewRequest("DELETE", h.urlPaymentID, h.requestBody)
	if err != nil {
		return ConfirmResponse{}, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return ConfirmResponse{}, err
	}

	return result, nil
}

func (h *CancelHandler) ReverseWithOrderID() (ConfirmResponse, error) {
	var (
		result ConfirmResponse
		url    = h.client.env.APIUrl
	)

	url += fmt.Sprintf("%s%s", "/1/physicalSales/orderId/", h.sale.MerchantOrderId)

	req, err := h.client.NewRequest("DELETE", url, h.requestBody)
	if err != nil {
		return ConfirmResponse{}, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return ConfirmResponse{}, err
	}

	return result, nil
}
