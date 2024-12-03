package main

import (
	"fmt"

	"github.com/silazemli/lab2-template/internal/services/payment"
)

func main() {
	db, err := payment.NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	srv := payment.NewServer(db)
	err = srv.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
