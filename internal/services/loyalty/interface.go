package loyalty

type loyaltyStorage interface {
	GetUser(username string) (Loyalty, error)
	IncrementCounter(username string) error
	DecrementCounter(username string) error
}
