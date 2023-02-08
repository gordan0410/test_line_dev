package app

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

type DBRepo interface {
	Get(req map[string]interface{}) (map[string]interface{}, error)
	GetAll(req map[string]interface{}, limit, offset int64) ([]map[string]interface{}, error)
	Create(req interface{}) (string, error)
	Close()
}

type MessageApp interface {
	Receive(events []*linebot.Event)
}

type messageApp struct {
	dbConn DBRepo
}

type DBsaveMsg struct {
	UserID string
	Msg    linebot.Message
}

func NewMessageApp(db DBRepo) MessageApp {
	return &messageApp{
		dbConn: db,
	}
}

func (r *messageApp) Receive(events []*linebot.Event) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			saveData := DBsaveMsg{
				UserID: event.Source.UserID,
				Msg:    event.Message,
			}
			r.dbConn.Create(saveData)
		}
	}
}
