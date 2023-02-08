/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgPath string

	rootCmd = &cobra.Command{
		Use:   "test_line_dev",
		Short: "test_line_dev is a line robot",
		Long:  `test_line_dev is a line robot, enjoy it!`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stdout, "please specify your command")
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgPath, "configPath", "c", "", "config file (default is $HOME/.test_line_dev.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(serverCmd)
}

func ErroHandle(err error) {
	fmt.Fprintln(os.Stderr, "sth wrong"+err.Error())
	os.Exit(1)
}
