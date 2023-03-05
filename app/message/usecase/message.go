package usecase

import (
	"test_line_dev/app/domain"
	"test_line_dev/app/tool"

	"github.com/line/line-bot-sdk-go/linebot"
)

type messageApp struct {
	dbConn domain.MessageRepo
	bot    *linebot.Client
}

func NewMessageApp(db domain.MessageRepo, bot *linebot.Client) domain.MessageApp {
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
			saveData := domain.CreateMessage{
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

func (r *messageApp) GetAllMsgByUserID(userID string, contentPerPage, page int) (res domain.GetAllMsgByUserIDRes, err error) {
	req := domain.GetMessage{
		UserID: userID,
	}

	allDataCounts, err := r.dbConn.CountDocuments(req)
	if err != nil {
		tool.ErroHandle("app", "MessageApp", "GetAllMsgByUserID", err)
		return res, err
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
			return res, err
		}

		result := make([]interface{}, len(datas))
		for i, v := range datas {
			result[i] = v["msg"]
		}
		return domain.GetAllMsgByUserIDRes{
			UserID:    userID,
			Messages:  result,
			TotalPage: totalPage,
			NowPage:   page,
		}, err
	}

	return res, err
}
