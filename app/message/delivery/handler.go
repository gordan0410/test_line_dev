package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"test_line_dev/tool"

	"github.com/gin-gonic/gin"
)

var contentPerPageMax = 10

type sendMsgReq struct {
	UserID string `bind:"required" json:"user_id"`
	Msg    string `bind:"required" json:"msg"`
}

func (s *Server) Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "123")
	}
}

func (s *Server) Receive() gin.HandlerFunc {
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

func (s *Server) Send() gin.HandlerFunc {
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

func (s *Server) GetAllMsgByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID != "" {
			var err error
			// default
			contentPerPage := 5
			nowPage := 1

			contentPerPageRaw := c.Query("content_per_page")
			if contentPerPageRaw != "" {
				contentPerPage, err = strconv.Atoi(contentPerPageRaw)
				if err != nil {
					tool.ErroHandle("server", "Server", "getAllMsgByUserID", err)
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				if contentPerPage > contentPerPageMax {
					contentPerPage = contentPerPageMax
				}
			}
			nowPageRaw := c.Query("page")
			if nowPageRaw != "" {
				nowPage, err = strconv.Atoi(nowPageRaw)
				if err != nil {
					tool.ErroHandle("server", "Server", "getAllMsgByUserID", err)
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
			}

			datas, err := s.messageApp.GetAllMsgByUserID(userID, contentPerPage, nowPage)
			if err != nil {
				tool.ErroHandle("server", "Server", "getAllMsgByUserID", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, datas)
		} else {
			err := errors.New("user_id not found")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
}
