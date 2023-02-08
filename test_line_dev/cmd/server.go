package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run robot server",
	Long:  `run robot server, enjoy it!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigType("json")
		fmt.Fprintln(os.Stdout, "robot server active")
		cfgPath, _ := cmd.Flags().GetString("configPath")
		if cfgPath != "" {
			file, err := ioutil.ReadFile(cfgPath)
			if err != nil {
				ErroHandle(err)
				return
			}

			viper.ReadConfig(bytes.NewBuffer(file))
		} else {
			viper.SetConfigName("config.json")
			wd, err := os.Getwd()
			if err != nil {
				ErroHandle(err)
				return
			}
			viper.AddConfigPath(wd + "/config/")
			viper.ReadInConfig()
		}
	},
}
