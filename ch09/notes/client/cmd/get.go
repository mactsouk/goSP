package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a specific note by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		resp, err := http.Get(fmt.Sprintf("%s/notes/%d", baseURL, id))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned: %s", resp.Status)
		}

		var note Note
		if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
			return err
		}

		fmt.Printf("[%d] %s — %s\n", note.ID, note.Title, note.Content)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
