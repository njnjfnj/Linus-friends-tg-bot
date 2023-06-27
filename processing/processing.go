package processing

import (
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"strconv"
	"time"

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

func (p *Processing) showMenu(chat_id int64, updates chan tgbotapi.Update, timer *time.Timer) {
	p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMenu))
	for {
		select {
		case <-timer.C:
			p.bot.Send(tgbotapi.NewMessage(chat_id, "MessageSessionTimeEnded"))
			return
		case upd := <-updates:
			if upd.Message != nil { // upd.FromChat().ChatConfig().ChatID == chat_id &&
				resetTimer(timer)
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
						if p.searchForProgrammers(chat_id, updates, user, timer) {
							return
						}
					case 2:
						if p.showProfileMenu(chat_id, updates, user, timer) {
							return
						}
					case 3:
						if p.showMatches(chat_id, updates, timer) {
							return
						}
					}
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "ERROR"))
				}

				p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMenu))
			}

		}
	}
}

func resetTimer(timer *time.Timer) {
	timer.Reset(20 * time.Second)
}
