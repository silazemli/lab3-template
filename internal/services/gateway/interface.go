package gateway

import (
	"github.com/silazemli/lab2-template/internal/services/loyalty"
	"github.com/silazemli/lab2-template/internal/services/reservation"
)

type gatewayAPI interface {
	GetData(username string) (string, error)
	GetUser(username string) (loyalty.Loyalty, error)
	GetAllHotels() ([]reservation.Hotel, error)
	GetAllReservations(username string) ([]reservation.Reservation, error)
	GetReservation(reservationUID string) (reservation.Reservation, error)
	MakeReservation(reservation reservation.Reservation) error
	CancelReservation(reservationUID string) error
}
