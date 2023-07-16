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
	bot                      *tgbotapi.BotAPI
	db                       storage.Storage
	adminPassword            string
	adminChoice              string
	timerResetDuration       int
	advertTimerResetDuration int
	advert                   *advertisement.Ad
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

func NewProcessing(botapi *tgbotapi.BotAPI, storage storage.Storage, adminPassword string, adminChoice string, adChan chan advertisement.Ad) *Processing {
	return &Processing{
		bot:                      botapi,
		db:                       storage,
		adminPassword:            adminPassword,
		adminChoice:              adminChoice,
		timerResetDuration:       30,
		advertTimerResetDuration: 1,
		advert:                   nil,
	}
}

func (p *Processing) showMenu(chat_id int64, updates chan tgbotapi.Update, timer *time.Timer, advertTimer *time.Timer) {
	for {
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMenu))
		select {
		case <-timer.C:
			return
		case <-advertTimer.C:
			if p.advert != nil {
				p.showAd(chat_id, p.advert, updates, timer)
			}
			p.resetAdvertTimer(advertTimer)
		case upd := <-updates:
			if upd.Message != nil {
				p.resetTimer(timer)

				user, err := p.db.GetUser(int(chat_id))
				if err != nil {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "Can not get user"+err.Error()))
					continue
				}
				switch upd.Message.Text {
				case "1":
					if p.searchForProgrammers(chat_id, updates, user, timer) { // , advertTimer
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

func (p *Processing) resetAdvertTimer(timer *time.Timer) {
	timer.Reset(time.Duration(p.advertTimerResetDuration) * time.Minute)
}

func (p *Processing) showAd(chat_id int64, MessageAdvert *advertisement.Ad, updates chan tgbotapi.Update, timer *time.Timer) (rating int) {
	defer func() {
		if rating != 0 {
			MessageAdvert.Rate += rating
			MessageAdvert.Rated++
		}
		MessageAdvert.Seen++
		p.db.UpdateAdRatingAndViews(MessageAdvert.Rate, MessageAdvert.Seen, MessageAdvert.Rated, MessageAdvert.Advert_id)
	}()
	p.processAdvertMessage(chat_id, MessageAdvert)

	var MessageErrorRating = "Enter a number from 0 to 5"

getRespondLoop1:
	for {
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageRateAdvert))
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
						break getRespondLoop1
					} else if bufRating > 0 && bufRating < 6 {
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

	return rating
}
