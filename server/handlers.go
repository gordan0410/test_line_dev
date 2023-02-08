package server

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func (s *Server) index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "123")
	}
}

func (s *Server) receiver() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := viper.GetString("channel_secret")
		token := viper.GetString("channel_access_token")
		bot, err := linebot.New(secret, token)
		if err != nil {
			log.Error().Caller().Err(err)
		}
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			log.Error().Caller().Err(err)
		}

		// app
		s.messageApp.Receive(events)

		// api

		c.String(200, "123")
	}
}
