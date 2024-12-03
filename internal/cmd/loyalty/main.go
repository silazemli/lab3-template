package main

import (
	"fmt"

	"github.com/silazemli/lab2-template/internal/services/loyalty"
)

func main() {
	db, err := loyalty.NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	srv := loyalty.NewServer(db)
	err = srv.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
