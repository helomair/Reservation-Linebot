package commandflow

import (
	"HPE-golang-test/services/component"
	"log"

	lineSDK "github.com/line/line-bot-sdk-go/v7/linebot"
)

func FlowControl(command string, params []string) lineSDK.SendingMessage {
	var ret lineSDK.SendingMessage

	log.Println(command)
	log.Println(params)

	switch command {
	case "?":
		ret = component.AvailableCommandsList()
	}

	return ret
}