package main

import (
	"fmt"
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

	// IDs to look up
	ids := []types.Uint128{
		types.ToUint128(1),
		types.ToUint128(2),
	}

	accounts, err := client.LookupAccounts(ids)
	if err != nil {
		log.Fatalf("Could not lookup accounts: %s", err)
	}

	fmt.Println("--- Account Balances ---")
	for _, acc := range accounts {
		fmt.Printf("Account ID: %s\n", acc.ID)
		fmt.Printf("  Ledger: %d\n", acc.Ledger)
		// Credits = Money received
		fmt.Printf("  Credits Posted: %s\n", acc.CreditsPosted)
		// Debits = Money sent
		fmt.Printf("  Debits Posted:  %s\n", acc.DebitsPosted)
		fmt.Println("------------------------")
	}
}
