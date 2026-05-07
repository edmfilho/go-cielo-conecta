package go_cielo_conecta

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
)

// GetParam defines the type for parameters used in GetPaymentBy method to specify the search criteria (PaymentId or MerchantOrderId).
type GetParam string

type Info struct {
	OrderID   string
	Amount    float64
	ProductID uint
}

var (
	PaymentID       = GetParam("PaymentId")
	MerchantOrderID = GetParam("MerchantOrderId")
)

// CreatePayment initializes a new payment with the provided order ID, amount, and product ID.
// It sets default values for installments, interest, capture, and payment date/time.
// The amount is converted to cents and rounding to the nearest integer.
//
// The method returns a SaleInterface that can be used to further customize the sale or execute it.
func (c *Client) CreatePayment(payment Info) SaleInterface {
	p := Payment{
		Installments:           1,          // Can be changed with SetInstallments().
		Interest:               ByMerchant, // Can be changed with SetInterest().
		Capture:                true,
		PaymentDateTime:        time.Now().Format("2006-01-02T15:04:05"),
		Amount:                 uint64(math.Round(payment.Amount * 100)),
		ProductId:              payment.ProductID,
		SubordinatedMerchantId: c.env.merchant.ID,
	}

	s := Sale{
		MerchantOrderId: payment.OrderID,
		Payment:         &p,
	}

	return newSaleHandler(c, s)
}

// GetPaymentBy retrieves a payment based on the specified parameter (PaymentId or MerchantOrderId) and query value.
// It constructs the appropriate endpoint URL based on the parameter and query, and optionally includes a transaction date.
// The method sends a GET requestBody to the API and returns the retrieved Sale object or an error if the requestBody fails.
//
// GET /1/physicalSales/{PaymentId}
// GET /1/physicalSales/MerchantOrderId/{MerchantOrderId}
func (c *Client) GetPaymentBy(param GetParam, query string, transactionDate ...time.Time) (sale *Sale, err error) {
	var (
		req      *http.Request
		endpoint = "/1/physicalSales"
	)

	switch param {
	case PaymentID:
		endpoint += fmt.Sprintf("/%s", query)
	case MerchantOrderID:
		endpoint += fmt.Sprintf("/MerchantOrderId/%s", query)
	default:
		return nil, errors.New("invalid param")
	}

	if len(transactionDate) > 0 {
		endpoint += fmt.Sprintf("?transactionDate=%s", transactionDate[0].Format("2006/01/02"))
	}

	req, err = c.NewRequest("GET", fmt.Sprintf("%s%s", c.env.APIQueryUrl, query), nil)
	if err != nil {
		return nil, err
	}

	err = c.Send(req, &sale)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (c *Client) ReversePayment(s Sale, issuerScriptsResults ...string) (ConfirmResponse, error) {
	cancel, err := newCancelHandler(c, s, issuerScriptsResults...)
	if err != nil {
		return ConfirmResponse{}, err
	}

	if s.Payment.PaymentId != "" {
		return cancel.ReverseWithPaymentID()
	}

	return cancel.ReverseWithOrderID()
}
