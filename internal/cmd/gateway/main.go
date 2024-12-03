package main

import (
	"fmt"

	"github.com/silazemli/lab2-template/internal/services/gateway"
)

func main() {
	srv := gateway.NewServer()
	err := srv.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
