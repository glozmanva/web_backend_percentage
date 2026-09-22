package main

import (
	"log"

	"deposit_month/internal/api"
)

func main() {
	log.Println("deposit_month application start")

	if err := api.StartServer(); err != nil {
		log.Fatal(err)
	}
}
