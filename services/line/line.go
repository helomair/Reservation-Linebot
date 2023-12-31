package line

import (
	"HPE-golang-test/configs"
	"HPE-golang-test/services/logger"
	"errors"
	"log"
	"net/http"
	"sync"

	lineSDK "github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineBotService struct {
	bot *lineSDK.Client
}

var mutex sync.Mutex
var lineBotService *LineBotService

func GetLineBotServiceInstance() *LineBotService {
	if lineBotService == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if lineBotService == nil {
			bot, err := lineSDK.New(configs.Configs.LineInfo.Secret, configs.Configs.LineInfo.AccessToken)
			if err != nil {
				panic("Line bot connection failed! in service/line/setup.go : " + err.Error())
			}

			lineBotService = &LineBotService{bot: bot}
		}
	}

	return lineBotService
}

func (handler *LineBotService) ParseRequestAndMakeMessage(request *http.Request) (LineEventMessageHandler, error) {
	events, err := handler.bot.ParseRequest(request)
	ret := LineEventMessageHandler{}
	if err != nil {
		log.Println(err.Error())
		return ret, err
	}

	if len(events) == 0 {
		return ret, errors.New("some error occurs when parse request")
	}

	ret.Event = events[0]

	return ret, nil
}

func (handler *LineBotService) Push(message LineEventMessageHandler) {
	target := message.UserId
	if target != "" {
		_, err := handler.bot.PushMessage(target, message.Message).Do()
		logger.ErrorFunc(err)
	}
}

func (handler *LineBotService) Reply(message LineEventMessageHandler) {
	replyToken := message.ReplyToken
	if replyToken != "" {
		_, err := handler.bot.ReplyMessage(replyToken, message.Message).Do()
		logger.ErrorFunc(err)
	}
}

func (handler *LineBotService) Broadcast(message string) error {
	_, err := handler.bot.BroadcastMessage(lineSDK.NewTextMessage(message)).Do()
	logger.ErrorFunc(err)
	return err
}
