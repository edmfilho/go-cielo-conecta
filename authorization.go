package go_cielo_conecta

import (
	"math"
	"time"
)

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

	return NewSaleHandler(c, &Sale{
		MerchantOrderId: orderId,
		Payment:         &p,
	})
}
