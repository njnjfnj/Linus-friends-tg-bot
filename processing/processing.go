package processing

import (
	"LinusFriends/advertisement"
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processing struct {
	bot                *tgbotapi.BotAPI
	db                 storage.Storage
	adminPassword      string
	adminChoice        string
	timerResetDuration int
	advert             chan advertisement.Ad
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

func NewProcessing(botapi *tgbotapi.BotAPI, storage storage.Storage, adminPassword *string, adminChoice *string, adChan chan advertisement.Ad) *Processing {
	return &Processing{
		bot:                botapi,
		db:                 storage,
		adminPassword:      *adminPassword,
		adminChoice:        *adminChoice,
		timerResetDuration: 30,
		advert:             adChan,
	}
}

func (p *Processing) showMenu(chat_id int64, updates chan tgbotapi.Update, timer *time.Timer) {
	for {
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMenu))
		select {
		case <-timer.C:
			return
		case MessageAdvert := <-p.advert:
			p.showAd(chat_id, &MessageAdvert, updates, timer)
		case upd := <-updates:
			if upd.Message != nil { // upd.FromChat().ChatConfig().ChatID == chat_id &&
				p.resetTimer(timer)

				user, err := p.db.GetUser(int(chat_id))
				if err != nil {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "Can not get user"+err.Error()))
					continue
				}
				switch upd.Message.Text {
				case "1":
					if p.searchForProgrammers(chat_id, updates, user, timer) {
						return
					}
				case "2":
					if p.showProfileMenu(chat_id, updates, user, timer) {
						return
					}
				case "3":
					if p.showMatches(chat_id, updates, timer) {
						return
					}
				case p.adminChoice:
					p.processAdmin(chat_id, updates, timer)
				}
			}
		}
	}
}

func (p *Processing) resetTimer(timer *time.Timer) {
	timer.Reset(time.Duration(p.timerResetDuration) * time.Minute)
}

func (p *Processing) showAd(chat_id int64, MessageAdvert *advertisement.Ad, updates chan tgbotapi.Update, timer *time.Timer) (rating int, seen int) {
	defer func() {

	}()
	p.bot.Send(MessageAdvert.Content)

	var MessageErrorRating = "Enter a number from 0 to 5"

getRespondLoop1:
	for {
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMenu))
		select {
		case <-timer.C:
			return
		case upd := <-updates:
			if upd.Message != nil {
				p.resetTimer(timer)

				if len(upd.Message.Text) == 1 {
					bufRating, err := strconv.Atoi(upd.Message.Text)
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, MessageErrorRating))
						continue
					}

					if bufRating == 0 {
						seen++
						break getRespondLoop1
					} else if bufRating > 0 && bufRating < 6 {
						seen++
						rating = bufRating
						break getRespondLoop1
					} else {
						p.bot.Send(tgbotapi.NewMessage(chat_id, MessageErrorRating))
					}
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, MessageErrorRating))
				}
			}
		}
	}

	MessageAdvert = &advertisement.Ad{}

	return rating, seen
}
