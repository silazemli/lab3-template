package reservation

type hotelStorage interface {
	GetAll() ([]Hotel, error)
	GetHotelID(hotelUID string) (int, error)
	GetHotel(ID string) (Hotel, error)
}

type reservationStorage interface {
	GetReservations(username string) ([]Reservation, error)
	GetReservation(reservationUID string) (Reservation, error)
	MakeReservation(reservation Reservation) error
	CancelReservation(reservationUID string) error
}
