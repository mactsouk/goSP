package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update an existing note by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ID: %v", err)
		}

		note := Note{Title: title, Content: content}
		data, _ := json.Marshal(note)

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/notes/%d", baseURL, id), bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned: %s", resp.Status)
		}

		fmt.Println("Note updated successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "New title for the note")
	updateCmd.Flags().StringVarP(&content, "content", "c", "", "New content for the note")
}
