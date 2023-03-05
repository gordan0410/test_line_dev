package domain

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

type MessageApp interface {
	GetBot() *linebot.Client
	Receive(events []*linebot.Event) error
	Send(userID, msg string) error
	GetAllMsgByUserID(userID string, contentPerPage, page int) (GetAllMsgByUserIDRes, error)
}

type MessageRepo interface {
	Get(req GetMessage) (map[string]interface{}, error)
	GetAll(req GetMessage, limit, offset int64) ([]map[string]interface{}, error)
	Create(req CreateMessage) (string, error)
	CountDocuments(req GetMessage) (int, error)
	Close()
}

type CreateMessage struct {
	UserID string
	Msg    linebot.Message
}
type GetMessage struct {
	UserID string
}

type GetAllMsgByUserIDRes struct {
	UserID    string        `json:"user_id"`
	Messages  []interface{} `json:"messages"`
	TotalPage int           `json:"total_page"`
	NowPage   int           `json:"now_page"`
}
