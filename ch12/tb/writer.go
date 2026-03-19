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

	// --- Step 1: Create Accounts ---
	accounts := []types.Account{
		{
			ID:     types.ToUint128(1),
			Ledger: 1,
			Code:   1,
		},
		{
			ID:     types.ToUint128(2),
			Ledger: 1,
			Code:   1,
		},
	}

	accountErrors, err := client.CreateAccounts(accounts)
	if err != nil {
		log.Fatalf("Fatal network error creating accounts: %s", err)
	}

	// TigerBeetle returns a list of errors for specific items in the batch.
	// We iterate through them to see if any specific account failed.
	for _, err := range accountErrors {
		// Example: checking if the account already exists is often useful logic
		if err.Result == types.AccountExists {
			log.Printf("Account at index %d already exists.", err.Index)
		} else {
			log.Printf("Error creating account at index %d: %s", err.Index, err.Result)
		}
	}
	
	// If no errors (or only 'exists' errors), we assume accounts are ready.
	fmt.Println("Accounts 1 and 2 are ready.")

	// --- Step 2: Create Transfer ---
	transfers := []types.Transfer{
		{
			ID:              types.ToUint128(100), // Unique ID for this transaction
			DebitAccountID:  types.ToUint128(1),
			CreditAccountID: types.ToUint128(2),
			Amount:          types.ToUint128(100),
			Ledger:          1,
			Code:            1,
		},
	}

	transferErrors, err := client.CreateTransfers(transfers)
	if err != nil {
		log.Fatalf("Fatal network error creating transfer: %s", err)
	}

	for _, err := range transferErrors {
		log.Printf("Error creating transfer at index %d: %s", err.Index, err.Result)
	}

	if len(transferErrors) == 0 {
		fmt.Println("Transfer of 100 units successful.")
	}
}
