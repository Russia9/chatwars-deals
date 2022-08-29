package app

import (
	"encoding/json"
	"gitea.russia9.dev/Russia9/chatwars-deals/messages"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"
)

type App struct {
	Bot      *telebot.Bot
	Chat     *telebot.Chat
	Consumer *kafka.Consumer
}

func Init(bot *telebot.Bot, chat *telebot.Chat, consumer *kafka.Consumer) error {
	app := App{
		Bot:      bot,
		Chat:     chat,
		Consumer: consumer,
	}

	err := consumer.SubscribeTopics([]string{"cw2-deals"}, nil)
	if err != nil {
		return err
	}
	channel := make(chan messages.DealMessage)
	go app.Sender(channel)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var message messages.DealMessage
			err = json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Error().Err(err).Str("module", "decoder").Send()
			}

			channel <- message

			log.Info().Str("partition", msg.TopicPartition.String()).Bytes("value", msg.Value).Send()
		} else {
			log.Error().Err(err).Str("module", "consumer").Send()
		}
	}
}
