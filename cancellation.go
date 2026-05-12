package go_cielo_conecta

import (
	"fmt"
	"strings"
)

type CancelInterface interface {
	ReverseWithPaymentID(paymentID string) (ConfirmResponse, error)
	ReverseWithOrderID(orderID string) (ConfirmResponse, error)
}

type CancelHandler struct {
	client       *Client
	requestBody  map[string]string
	urlPaymentID string
	urlOrderID   string
}

func newCancelHandler(c *Client, emvData string, issuerScriptsResults ...string) CancelInterface {
	urlPaymentID := fmt.Sprintf("%s/1/physicalSales/{PaymentID}", c.env.APIUrl)
	urlOrderID := fmt.Sprintf("%s/1/physicalSales/orderId/{OrderID}", c.env.APIUrl)

	body := map[string]string{
		"EmvData":             emvData,
		"IssuerScriptResults": "0000",
	}

	if len(issuerScriptsResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptsResults[0]
	}

	return &CancelHandler{client: c, requestBody: body, urlPaymentID: urlPaymentID, urlOrderID: urlOrderID}
}

func (h *CancelHandler) ReverseWithPaymentID(paymentID string) (ConfirmResponse, error) {
	var result ConfirmResponse

	req, err := h.client.NewRequest("DELETE", strings.Replace(h.urlPaymentID, "{PaymentID}", paymentID, 1), h.requestBody)
	if err != nil {
		return ConfirmResponse{}, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return ConfirmResponse{}, err
	}

	return result, nil
}

func (h *CancelHandler) ReverseWithOrderID(orderID string) (ConfirmResponse, error) {
	var result ConfirmResponse

	req, err := h.client.NewRequest("DELETE", strings.Replace(h.urlOrderID, "{OrderID}", orderID, 1), h.requestBody)
	if err != nil {
		return ConfirmResponse{}, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return ConfirmResponse{}, err
	}

	return result, nil
}
