package go_cielo_conecta

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// CreatePayment initializes a new payment with the provided order ID, amount, and product ID.
// It sets default values for installments, interest, capture, and payment date/time.
// The amount is converted to cents and rounding to the nearest integer.
//
// The method returns a SaleInterface that can be used to further customize the sale or execute it.
func (c *Client) CreatePayment(orderId string, amount float64, productId uint) SaleInterface {
	p := Payment{
		Installments:           1,
		Interest:               ByMerchant, // Initialized with ByMerchant, but can be changed with SetInterest().
		Capture:                true,
		PaymentDateTime:        time.Now().Format("2006-01-02T15:04:05"),
		Amount:                 uint64(math.Round(amount * 100)),
		ProductId:              productId,
		SubordinatedMerchantId: c.merchant.ID,
	}

	return newSaleHandler(c, &Sale{
		MerchantOrderId: orderId,
		Payment:         &p,
	})
}

// GetPaymentByID retrieves a payment by its unique identifier (PaymentId).
// It constructs the appropriate API query URL and sends a GET request to the Cielo Conecta API.
//
// GET {APIQueryUrl}/1/physicalSales/{PaymentId}
func (c *Client) GetPaymentByID(paymentID string, transactionDate ...time.Time) (sale *Sale, err error) {
	var (
		query = fmt.Sprintf("/1/physicalSales/%s", paymentID)
		req   *http.Request
	)

	if len(transactionDate) > 0 {
		query += fmt.Sprintf("?transactionDate=%s", transactionDate[0].Format("2006/01/02"))
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

// GET {APIQueryUrl}/1/physicalSales/MerchantOrderId/{MerchantOrderId}
func (c *Client) GetPaymentByMerchantOrderID(merchantOrderID string, transactionDate ...time.Time) (sale *Sale, err error) {
	var (
		query = fmt.Sprintf("/1/physicalSales/MerchantOrderId/%s", merchantOrderID)
		req   *http.Request
	)

	if len(transactionDate) > 0 {
		query += fmt.Sprintf("?transactionDate=%s", transactionDate[0].Format("2006/01/02"))
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
