package loyalty

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type storage struct {
	db *gorm.DB
}

func NewDB() (*storage, error) {
	dsn := os.Getenv("LOYALTY_DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &storage{}, err
	}
	return &storage{db}, nil
}

func (stg *storage) GetUser(username string) (Loyalty, error) {
	loyalty := Loyalty{}
	err := stg.db.Table("loyalty").Where("username = ?", username).Take(&loyalty).Error
	if err != nil {
		return Loyalty{}, err
	}
	return loyalty, nil
}

func (stg *storage) IncrementCounter(username string) error {
	loyalty := Loyalty{}
	err := stg.db.Table("loyalty").Where("username = ?", username).Take(&loyalty).Error
	if err != nil {
		return err
	}
	loyalty.ReservationCount += 1
	UpdateStatus(&loyalty)
	err = stg.db.Table("loyalty").Where("username = ?", username).Updates(&loyalty).Error
	if err != nil {
		return err
	}
	return nil
}

func (stg *storage) DecrementCounter(username string) error {
	loyalty := Loyalty{}
	err := stg.db.Table("loyalty").Where("username = ?", username).Take(&loyalty).Error
	if err != nil {
		return err
	}
	loyalty.ReservationCount -= 1
	UpdateStatus(&loyalty)
	err = stg.db.Table("loyalty").Where("username = ?", username).Updates(&loyalty).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateStatus(loyalty *Loyalty) {
	numberOfReservations := loyalty.ReservationCount
	if numberOfReservations >= 20 {
		loyalty.Status = "GOLD"
		loyalty.Discount = 10
	} else if numberOfReservations >= 10 {
		loyalty.Status = "SILVER"
		loyalty.Discount = 7
	} else {
		loyalty.Status = "BRONZE"
		loyalty.Discount = 5
	}
}
