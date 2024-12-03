package gateway

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LoyaltyService     string `env:"LOYALTY_SERVICE"`
	PaymentService     string `env:"PAYMENT_SERVICE"`
	ReservationService string `env:"RESERVATION_SERVICE"`
}

func NewConfig() *Config {
	localPath := "./configs/gateway.env"
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil
	}

	err = cleanenv.ReadConfig(localPath, &cfg)
	if err != nil {
		return nil
	}

	return &cfg
}
