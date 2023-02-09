package server

import (
	"context"
	"fmt"
	"os"

	"test_line_dev/app"

	"github.com/gin-gonic/gin"
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

	f, err := os.Create("urls.txt")
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("linebot webhook url:\n%s\nlinebot api:\n%s\n", tun.URL()+"/linebot", tun.URL()+"/linebot/message")
	_, err = f.Write([]byte(msg))
	f.Close()

	r := s.Rounters()
	err = r.RunListener(tun)
	if err != nil {
		return err
	}
	return nil
}
