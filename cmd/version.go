package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of csv_json",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CSV2JSON v0.1 -- HEAD")
	},
}
