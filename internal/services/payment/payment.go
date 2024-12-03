package payment

type Payment struct {
	PaymentUID string `json:"paymentUid"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}
