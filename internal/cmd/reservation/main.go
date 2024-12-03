package main

import (
	"fmt"

	"github.com/silazemli/lab2-template/internal/services/reservation"
)

func main() {
	hdb, err := reservation.NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	rdb, err := reservation.NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	srv := reservation.NewServer(rdb, hdb)
	err = srv.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
