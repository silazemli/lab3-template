package loyalty

type Loyalty struct {
	Username         string `json:"username"`
	ReservationCount int    `json:"reservationCount"`
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
}
