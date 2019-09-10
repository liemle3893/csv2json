package cmd

import (
	c "encoding/config"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "csv2json",
	Short: "Convert CSV into JSON",
	Long:  `A Fast and Flexible Converter to Support convert CSV to JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		configFile, _ := cmd.Flags().GetString("config-file")
		f, err := os.Open(configFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		configText, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		config, err := c.ParseConfig(string(configText))
		if err != nil {
			fmt.Println(err)
			os.Exit(100)
		}
		config.Exec()
	},
}

// Execute run csv2json
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("config-file", "c", "config.hcl", "Configuration file path")
}
