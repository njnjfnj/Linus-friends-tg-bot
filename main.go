package main

import (
	"context"
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
	if err := db.Init(context.TODO()); err != nil {
		log.Fatal("Can not init a db: ", err)
	}

	// bot init
	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	UsersInProcess := make(map[int]chan tgbotapi.Update)

	// creating processing object
	// cheking updates
	updates := bot.GetUpdatesChan(u)
	process := processing.NewProcessing(bot, db)
	for upd := range updates {
		go func(update tgbotapi.Update) {
			if update.Message != nil {
				chat_id := update.FromChat().ChatConfig().ChatID
				if UsersInProcess[int(chat_id)] == nil {
					if len(update.Message.Text) != 0 && update.Message.Text[0] == '/' {
						UsersInProcess[int(chat_id)] = make(chan tgbotapi.Update)
						if err := process.CMD(strings.ReplaceAll(update.Message.Text, " ", ""), chat_id, UsersInProcess[int(chat_id)]); err != nil {
							log.Println("Can not process update: ", err)
						}
						close(UsersInProcess[int(chat_id)])
						UsersInProcess[int(chat_id)] = nil
					} else {
						bot.Send(tgbotapi.NewMessage(chat_id, "Unknown message"))
						return
					}
				} else {
					UsersInProcess[int(chat_id)] <- upd
				}
			}
		}(upd)
	}
}
