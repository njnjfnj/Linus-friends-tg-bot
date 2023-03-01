package main

import (
	"flag"
	"log"
	"strings"

	"LinusFriends/processing"
	Database "LinusFriends/storage/sqlite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	//getting token
	token := flag.String(
		"token",
		"",
		"token for access to tg bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	// creating db
	db, err := Database.NewDatabase("db/storage.db")
	if err != nil {
		log.Fatal("Can not make a db: ", err)
	}
	// if err := db.Init(context.TODO()); err != nil { когда будешь работать с бд тогда и раскоменть
	// 	log.Fatal("Can not init a db: ", err)
	// }

	// bot init
	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// creating processing object
	process := processing.NewProcessing(bot, db)

	// cheking updates
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			var err error
			chat_id := update.FromChat().ChatConfig().ChatID
			if len(update.Message.Text) != 0 {
				if update.Message.Text[0] == '/' {
					err = process.CMD(strings.ReplaceAll(update.Message.Text, " ", ""), chat_id)
				} else {
					err = process.Text(update.Message.Text, chat_id)
				}
			} else if update.Message.Photo != nil {
				err = process.Photo(update.Message.Photo, chat_id)
			} else if _, err := bot.Send(tgbotapi.NewMessage(chat_id, "Unknown message")); err != nil {
				log.Print("Can not send message to the ", chat_id, ": ", err)
				continue
			}

			if err != nil {
				log.Println("Can not process update: ", err)
			}
		}
	}
}
