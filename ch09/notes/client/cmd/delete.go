package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a note by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/notes/%d", baseURL, id), nil)
		if err != nil {
			return err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("server returned: %s", resp.Status)
		}

		fmt.Println("Note deleted successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
