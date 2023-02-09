package server

import "github.com/gin-gonic/gin"

func (s *Server) Rounters() *gin.Engine {
	router := s.router
	router.GET("/", s.index())
	router.POST("/linebot", s.receive())
	router.POST("/linebot/message", s.send())
	router.GET("/linebot/message", s.getAllMsgByUserID())

	return router
}
