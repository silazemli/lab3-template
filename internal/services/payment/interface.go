package payment

type paymentStorage interface {
	GetPayment(paymentUID string) (Payment, error)
	PostPayment(thePayment Payment) error
	CancelPayment(paymentUID string) error
}
