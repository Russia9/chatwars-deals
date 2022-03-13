package app

import (
	"cw-deals-watcher/messages"
	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
	"time"
)

func (a *App) Sender(channel chan messages.DealMessage) {
	messagePool := make(map[string][]messages.DealMessage)
	lastTimestamp := int32(time.Now().Unix())
	counter := 0

	for {
		var message messages.DealMessage
		message = <-channel
		messagePool[message.Item] = append(messagePool[message.Item], message)
		counter++

		if int32(time.Now().Unix())-20 >= lastTimestamp || counter >= 30 {
			msgString := ""
			for item, value := range messagePool {
				msgString += "<b>" + item + "</b>:\n"
				for _, msg := range value {
					msgString += "<code>" + msg.SellerCastle + msg.SellerName + " -> " + msg.BuyerCastle + msg.BuyerName + ", " +
						"" + "" + strconv.Itoa(msg.Quantity) + " x " + strconv.Itoa(msg.Price) + "💰 </code>\n"
				}
			}

			_, err := a.Bot.Send(a.Chat, msgString, telebot.ModeHTML)
			if err != nil {
				log.Error().Err(err).Str("module", "sender").Send()
			}

			lastTimestamp = int32(time.Now().Unix())
			counter = 0
			messagePool = make(map[string][]messages.DealMessage)
		}
	}
}
