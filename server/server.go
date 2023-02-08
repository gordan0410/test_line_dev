package server

import (
	"context"

	"test_line_dev/app"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

type Server struct {
	router     *gin.Engine
	messageApp app.MessageApp
}

func NewServer(r *gin.Engine, receiver app.MessageApp) *Server {
	return &Server{
		router:     r,
		messageApp: receiver,
	}
}

func (s *Server) Run() error {
	tun, err := ngrok.Listen(context.Background(),
		config.HTTPEndpoint(),
	)
	if err != nil {
		return err
	}
	log.Info().Msgf("libot webhook url : %s", tun.URL()+"/linebot")
	r := s.Rounters()
	err = r.RunListener(tun)
	if err != nil {
		return err
	}
	return nil
}
