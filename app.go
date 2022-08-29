package main

import (
	"gitea.russia9.dev/Russia9/chatwars-deals/app"
	"gitea.russia9.dev/Russia9/chatwars-deals/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"
	"os"
	"strconv"
	"time"
)

func main() {
	// Log settings
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	pretty, err := strconv.ParseBool(os.Getenv("LOG_PRETTY"))
	if err != nil {
		pretty = false
	}
	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	switch os.Getenv("LOG_LEVEL") {
	case "DISABLED":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Kafka consumer init
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": utils.GetEnv("KAFKA_ADDRESS", "localhost"),
		"group.id":          "cw2-deals",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		log.Fatal().Err(err).Str("module", "kafka").Send()
	}

	// Bot init
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal().Err(err).Str("module", "bot").Send()
	}

	chat, err := bot.ChatByID(utils.GetEnv("TELEGRAM_CHANNEL", "-1001429147053"))
	if err != nil {
		log.Fatal().Err(err).Str("module", "bot").Send()
	}

	err = app.Init(bot, chat, consumer)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
