package go_cielo_conecta

import (
	"errors"
	"fmt"
	"net/http"
)

type CancelInterface interface {
	UndoWithPaymentID(issuerScriptsResults ...string) (*ConfirmResponse, error)
	UndoWithMerchantOrderID(issuerScriptsResults ...string) (*ConfirmResponse, error)
}

type CancelHandler struct {
	client *Client
	sale   *Sale
}

func newCancelHandler(client *Client, sale Sale) (CancelInterface, error) {
	if sale.Payment == nil {
		return nil, errors.New("payment is required")
	}

	if sale.Payment.CreditCard == nil {
		return nil, errors.New("credit_card is required")
	}

	ss := &Sale{
		MerchantOrderId: sale.MerchantOrderId,
		Payment:         sale.Payment,
	}

	return &CancelHandler{
		client: client,
		sale:   ss,
	}, nil
}

func (h *CancelHandler) UndoWithPaymentID(issuerScriptsResults ...string) (result *ConfirmResponse, err error) {
	var (
		body = map[string]string{}
		url  = h.client.env.APIUrl
		req  *http.Request
	)

	body["EmvData"] = h.sale.Payment.CreditCard.EmvData
	body["IssuerScriptResults"] = "0000"

	if len(issuerScriptsResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptsResults[0]
	}

	url += fmt.Sprintf("%s%s", "/1/physicalSales/", h.sale.Payment.PaymentId)

	req, err = h.client.NewRequest("DELETE", url, body)
	if err != nil {
		return nil, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *CancelHandler) UndoWithMerchantOrderID(issuerScriptsResults ...string) (result *ConfirmResponse, err error) {
	var (
		body = map[string]string{}
		url  = h.client.env.APIUrl
		req  *http.Request
	)

	body["EmvData"] = h.sale.Payment.CreditCard.EmvData
	body["IssuerScriptResults"] = "0000"

	if len(issuerScriptsResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptsResults[0]
	}

	url += fmt.Sprintf("%s%s", "/1/physicalSales/orderId/", h.sale.MerchantOrderId)

	req, err = h.client.NewRequest("DELETE", url, body)
	if err != nil {
		return nil, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
