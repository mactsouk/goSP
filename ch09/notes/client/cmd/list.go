package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get(baseURL + "/notes")
		if err != nil {
			return fmt.Errorf("error fetching notes: %w", err)
		}
		defer resp.Body.Close()

		var notes []Note
		if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return nil
		}

		for _, n := range notes {
			fmt.Printf("[%d] %s — %s\n", n.ID, n.Title, n.Content)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
