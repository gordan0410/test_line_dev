package server

import "github.com/gin-gonic/gin"

func (s *Server) Rounters() *gin.Engine {
	router := s.router
	router.GET("/", s.index())
	router.POST("/linebot", s.receiver())
	router.POST("/linebot/send", s.sendMsg())
	return router
}
