package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new note",
	RunE: func(cmd *cobra.Command, args []string) error {
		note := Note{Title: title, Content: content}
		data, _ := json.Marshal(note)

		resp, err := http.Post(baseURL+"/notes", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("error creating note: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("server returned: %s", resp.Status)
		}

		fmt.Println("Note created successfully!")
		return nil
	},
}

var title, content string

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&title, "title", "t", "", "Note title")
	addCmd.Flags().StringVarP(&content, "content", "c", "", "Note content")
	addCmd.MarkFlagRequired("title")
	addCmd.MarkFlagRequired("content")
}
