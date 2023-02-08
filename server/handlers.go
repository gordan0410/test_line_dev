package server

import (
	"net/http"

	"test_line_dev/tool"

	"github.com/gin-gonic/gin"
)

type sendMsgReq struct {
	UserID string `bind:"required" json:"user_id"`
	Msg    string `bind:"required" json:"msg"`
}

func (s *Server) index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "123")
	}
}

func (s *Server) receiver() gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := s.messageApp.GetBot().ParseRequest(c.Request)
		if err != nil {
			tool.ErroHandle("server", "Server", "receiver", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// app
		err = s.messageApp.Receive(events)
		if err != nil {
			tool.ErroHandle("server", "Server", "receiver", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.String(200, "OK")
	}
}

func (s *Server) sendMsg() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req sendMsgReq
		if err := c.ShouldBindJSON(&req); err != nil {
			tool.ErroHandle("server", "Server", "sendMsg", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := s.messageApp.Send(req.UserID, req.Msg)
		if err != nil {
			tool.ErroHandle("server", "Server", "sendMsg", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, req)
	}
}
