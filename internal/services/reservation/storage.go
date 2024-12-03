package reservation

import (
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type storage struct {
	db *gorm.DB
}

func NewDB() (*storage, error) {
	dsn := os.Getenv("RESERVATION_DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &storage{}, err
	}
	return &storage{db}, nil
}

func (stg *storage) GetAll() ([]Hotel, error) {
	hotels := []Hotel{}
	err := stg.db.Table("hotels").Find(&hotels).Error
	if err != nil {
		log.Info().Msg("Failed at db level")
		return []Hotel{}, err
	}
	return hotels, nil
}

func (stg *storage) GetReservations(username string) ([]Reservation, error) {
	reservations := []Reservation{}
	err := stg.db.Table("reservation").Where("username = ?", username).Find(&reservations).Error
	if err != nil {
		return []Reservation{}, err
	}
	return reservations, nil
}

func (stg *storage) GetReservation(reservationUID string) (Reservation, error) {
	reservation := Reservation{}
	err := stg.db.Table("reservation").Where("reservation_uid = ?", reservationUID).Take(&reservation).Error
	if err != nil {
		return Reservation{}, err
	}
	return reservation, err
}

func (stg *storage) MakeReservation(reservation Reservation) error {
	err := stg.db.Table("reservation").Create(&reservation).Error
	if err != nil {
		return err
	}
	return nil
}

func (stg *storage) CancelReservation(reservationUID string) error {
	reservation := Reservation{}
	err := stg.db.Table("reservation").Where("reservation_uid = ?", reservationUID).Take(&reservation).Error
	if err != nil {
		return err
	}
	reservation.Status = "CANCELED"
	err = stg.db.Table("reservation").Where("reservation_uid = ?", reservationUID).Updates(&reservation).Error
	if err != nil {
		return err
	}
	return nil
}

func (stg *storage) GetHotelID(hotelUID string) (int, error) {
	var ID int
	err := stg.db.Table("hotels").Where("hotel_uid = ?", hotelUID).Select("id").Take(&ID).Error
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func (stg *storage) GetHotel(ID string) (Hotel, error) {
	var hotel Hotel
	err := stg.db.Table("hotels").Where("id = ?", ID).Take(&hotel).Error
	if err != nil {
		return Hotel{}, err
	}
	return hotel, nil
}
