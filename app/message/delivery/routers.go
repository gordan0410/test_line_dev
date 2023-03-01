package delivery

import "github.com/gin-gonic/gin"

func (s *Server) Rounters() *gin.Engine {
	router := s.router
	router.GET("/", s.Index())
	router.POST("/linebot", s.Receive())
	router.POST("/linebot/message", s.Send())
	router.GET("/linebot/message", s.GetAllMsgByUserID())

	return router
}
