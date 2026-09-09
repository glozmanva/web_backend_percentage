package main

import (
	"log"

	"web_backend_percentage/internal/api"
)

func main() {
	log.Println("Application start!")

	if err := api.StartServer(); err != nil {
		log.Fatal(err)
	}
}
