package go_cielo_conecta

import (
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
	urlOrderID   string
}

func newCancelHandler(c *Client, authorizedSale Sale, issuerScriptsResults ...string) (CancelInterface, error) {
	if authorizedSale.Payment == nil {
		return nil, ErrPaymentRequired
	}

	withPaymentID := fmt.Sprintf("%s/1/physicalSales/%s", c.env.APIUrl, authorizedSale.Payment.PaymentId)
	withOrderID := fmt.Sprintf("%s/1/physicalSales/orderId/%s", c.env.APIUrl, authorizedSale.MerchantOrderId)

	body := map[string]string{
		"EmvData":             authorizedSale.Payment.getEmvData(),
		"IssuerScriptResults": "0000",
	}

	if len(issuerScriptsResults) > 0 {
		body["IssuerScriptResults"] = issuerScriptsResults[0]
	}

	return &CancelHandler{client: c, sale: authorizedSale, requestBody: body, urlPaymentID: withPaymentID, urlOrderID: withOrderID}, nil
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
	var result ConfirmResponse

	req, err := h.client.NewRequest("DELETE", h.urlOrderID, h.requestBody)
	if err != nil {
		return ConfirmResponse{}, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return ConfirmResponse{}, err
	}

	return result, nil
}
