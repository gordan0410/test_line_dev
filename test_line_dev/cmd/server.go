package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"test_line_dev/app"
	"test_line_dev/repository"
	"test_line_dev/server"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
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
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dataBase := viper.GetString("mongoDB_database")
		collection := viper.GetString("mongoDB_collection")
		dbRepo := repository.NewMogoDB(ctx, dataBase, collection)
		defer dbRepo.Close()

		secret := viper.GetString("channel_secret")
		token := viper.GetString("channel_access_token")
		bot, err := linebot.New(secret, token)
		if err != nil {
			ErroHandle(err)
			return
		}
		receiverApp := app.NewMessageApp(dbRepo, bot)
		router := gin.Default()
		server := server.NewServer(router, receiverApp)
		server.Run()
	},
}
