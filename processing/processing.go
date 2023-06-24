package processing

import (
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processing struct {
	bot *tgbotapi.BotAPI
	db  storage.Storage
}

func (p *Processing) CMD(cmd string, chat_id int64, updates chan tgbotapi.Update) (err error) {
	defer func() { err = e.WrapIfErr("Can not process a cmd: ", err) }()

	switch cmd {
	case "/help":
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageHelp))
	case "/start":
		if err := p.processCMDStart(chat_id, updates); err != nil {
			return err
		}
	}
	return nil
}

func NewProcessing(botapi *tgbotapi.BotAPI, storage storage.Storage) *Processing {
	return &Processing{
		bot: botapi,
		db:  storage,
	}
}

func (p *Processing) showMenu(chat_id int64, updates chan tgbotapi.Update) {
	p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMenu))
	for upd := range updates {
		if upd.Message != nil { // upd.FromChat().ChatConfig().ChatID == chat_id &&
			if len(upd.Message.Text) == 1 {
				buf, err := strconv.ParseInt(upd.Message.Text, 0, 8)
				if err != nil {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "ERROR"))
				}
				check := int8(buf)
				user, err := p.db.GetUser(int(chat_id))
				if err != nil {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "Can not get user"+err.Error()))
					continue
				}
				switch check {
				case 1:
					p.searchForProgrammers(chat_id, updates, user)
				case 2:
					p.showProfileMenu(chat_id, updates, user)

				case 3:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "SHOWING..."))
				}
			} else {
				p.bot.Send(tgbotapi.NewMessage(chat_id, "ERROR"))
			}

			p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMenu))
		}
	}
}
