package main

import (
	"log"
	"os"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func main() {
	client, err := tb.NewClient(types.ToUint128(0), []string{"127.0.0.1:3000"})
	if err != nil {
		log.Printf("Error creating client: %s", err)
		os.Exit(1)
	}
	defer client.Close()

	log.Println("Successfully connected to TigerBeetle!")
}
