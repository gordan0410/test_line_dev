package app

import (
	"test_line_dev/tool"

	"github.com/line/line-bot-sdk-go/linebot"
)

type DBRepo interface {
	Get(req map[string]interface{}) (map[string]interface{}, error)
	GetAll(req map[string]interface{}, limit, offset int64) ([]map[string]interface{}, error)
	Create(req interface{}) (string, error)
	Close()
}

type MessageApp interface {
	GetBot() *linebot.Client
	Receive(events []*linebot.Event) error
	Send(userID, msg string) error
}

type messageApp struct {
	dbConn DBRepo
	bot    *linebot.Client
}

type DBsaveMsg struct {
	UserID string
	Msg    linebot.Message
}

func NewMessageApp(db DBRepo, bot *linebot.Client) MessageApp {
	return &messageApp{
		dbConn: db,
		bot:    bot,
	}
}

func (r *messageApp) GetBot() *linebot.Client {
	return r.bot
}

func (r *messageApp) Receive(events []*linebot.Event) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			saveData := DBsaveMsg{
				UserID: event.Source.UserID,
				Msg:    event.Message,
			}
			_, err := r.dbConn.Create(saveData)
			if err != nil {
				tool.ErroHandle("app", "MessageApp", "Receive", err)
				return err
			}
		}
	}
	return nil
}

func (r *messageApp) Send(userID, msg string) error {
	_, err := r.bot.PushMessage(userID, linebot.NewTextMessage(msg)).Do()
	if err != nil {
		return err
	}
	return nil
}
