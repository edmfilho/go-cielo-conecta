package go_cielo_conecta

import (
	"context"
	"fmt"
	"net/http"
)

type CancelInterface interface {
	TryReversePayment() (ConfirmResponse, error)
}

type CancelHandler struct {
	client       *Client
	ctx          context.Context
	data         ReverseRequest
	hasPaymentID bool
}

func newCancelHandler(ctx context.Context, c *Client, request ReverseRequest) CancelInterface {
	var (
		hasPaymentID = false
	)

	if request.PaymentID != "" {
		hasPaymentID = true
	}

	return &CancelHandler{
		client:       c,
		ctx:          ctx,
		data:         request,
		hasPaymentID: hasPaymentID,
	}
}

func (h *CancelHandler) TryReversePayment() (ConfirmResponse, error) {
	var (
		result ConfirmResponse
		req    *http.Request
		err    error
	)

	body := map[string]string{"EmvData": h.data.EmvData}

	if h.hasPaymentID {
		req, err = h.client.NewRequestWithContext(h.ctx, http.MethodDelete,
			fmt.Sprintf("%s/1/physicalSales/%s", h.client.env.APIUrl, h.data.PaymentID),
			body,
		)
	} else {
		req, err = h.client.NewRequestWithContext(h.ctx, http.MethodDelete,
			fmt.Sprintf("%s/1/physicalSales/MerchantOrderId/%s", h.client.env.APIUrl, h.data.MerchantOrderId),
			body,
		)
	}

	if err != nil {
		return ConfirmResponse{}, err
	}

	err = h.client.Send(req, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
