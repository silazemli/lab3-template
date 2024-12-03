package payment

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type storage struct {
	db *gorm.DB
}

func NewDB() (*storage, error) {
	dsn := os.Getenv("PAYMENT_DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &storage{}, err
	}
	return &storage{db}, nil
}

func (stg *storage) PostPayment(thePayment Payment) error {
	err := stg.db.Table("payment").Create(&thePayment).Error
	if err != nil {
		return err
	}
	return nil
}

func (stg *storage) GetPayment(paymentUID string) (Payment, error) {
	payment := Payment{}
	err := stg.db.Table("payment").Where("payment_uid = ?", paymentUID).Take(&payment).Error
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (stg *storage) CancelPayment(paymentUID string) error {
	payment := Payment{}
	err := stg.db.Table("payment").Where("payment_uid = ?", paymentUID).Take(&payment).Error
	if err != nil {
		return err
	}
	payment.Status = "CANCELED"
	err = stg.db.Table("payment").Where("payment_uid = ?", paymentUID).Updates(&payment).Error
	if err != nil {
		return err
	}
	return nil
}
