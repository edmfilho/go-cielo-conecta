package go_cielo_conecta

import (
	"math"
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
