package main

import (
	"log"

	"github.com/MurilojrMarques/mini-ecommerce-go/cmd/api"
)

func main() {
	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal("Falha ao iniciar o servidor:", err)
	}
}
