package app

import (
	"encoding/json"

	"test_line_dev/tool"

	"github.com/line/line-bot-sdk-go/linebot"
)

type DBRepo interface {
	Get(req map[string]interface{}) (map[string]interface{}, error)
	GetAll(req map[string]interface{}, limit, offset int64) ([]map[string]interface{}, error)
	Create(req interface{}) (string, error)
	CountDocuments(req map[string]interface{}) (int, error)
	Close()
}

type MessageApp interface {
	GetBot() *linebot.Client
	Receive(events []*linebot.Event) error
	Send(userID, msg string) error
	GetAllMsgByUserID(userID string, contentPerPage, page int) (GetAllMsgByUserIDRes, error)
}

type messageApp struct {
	dbConn DBRepo
	bot    *linebot.Client
}

type DBsaveMsg struct {
	UserID string
	Msg    linebot.Message
}

type GetAllMsgByUserIDRes struct {
	UserID    string   `json:"user_id"`
	Messages  []string `json:"messages"`
	TotalPage int      `json:"total_page"`
	NowPage   int      `json:"now_page"`
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
		tool.ErroHandle("app", "MessageApp", "Send", err)
		return err
	}
	return nil
}

func (r *messageApp) GetAllMsgByUserID(userID string, contentPerPage, page int) (GetAllMsgByUserIDRes, error) {
	req := map[string]interface{}{
		"userid": userID,
	}

	allDataCounts, err := r.dbConn.CountDocuments(req)
	if err != nil {
		tool.ErroHandle("app", "MessageApp", "GetAllMsgByUserID", err)
		return GetAllMsgByUserIDRes{}, err
	}

	totalPage := allDataCounts / contentPerPage
	rp := allDataCounts % contentPerPage
	if rp > 0 {
		totalPage += 1
	}

	if totalPage > 0 {
		offset := contentPerPage * (page - 1)
		datas, err := r.dbConn.GetAll(req, int64(contentPerPage), int64(offset))
		if err != nil {
			tool.ErroHandle("app", "MessageApp", "GetAllMsgByUserID", err)
			return GetAllMsgByUserIDRes{}, err
		}

		var result []string
		for _, v := range datas {
			data, err := json.Marshal(v["msg"])
			if err != nil {
				tool.ErroHandle("app", "MessageApp", "GetAllMsgByUserID", err)
				return GetAllMsgByUserIDRes{}, err
			}
			result = append(result, string(data))
		}
		return GetAllMsgByUserIDRes{
			UserID:    userID,
			Messages:  result,
			TotalPage: totalPage,
			NowPage:   page,
		}, nil
	}

	return GetAllMsgByUserIDRes{}, err
}
