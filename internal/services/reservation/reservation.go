package reservation

type Reservation struct {
	ReservationUID string `json:"reservation_uid"`
	Username       string `json:"username"`
	PaymentUID     string `json:"payment_uid"`
	HotelID        int    `json:"hotel_id"`
	Status         string `json:"status"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}
