package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var baseURL string

var rootCmd = &cobra.Command{
	Use:   "notescli",
	Short: "A CLI client for the Go note-taking web service",
	Long:  `notescli communicates with a running note-taking HTTP server to create, list, update, and delete notes.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "http://localhost:8080", "Base URL of the note server")
}
