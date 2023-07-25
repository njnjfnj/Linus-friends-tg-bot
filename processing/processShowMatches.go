package processing

import (
	"LinusFriends/storage"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) showMatches(chat_id int64, updates chan tgbotapi.Update, timer *time.Timer, advertTimer *time.Timer) bool {
	isAdvertTimer := false

	matches, err := p.db.GetMatches(chat_id)
	if len(matches) == 0 || err == storage.ErrNoFriends {
		p.bot.Send(tgbotapi.NewMessage(chat_id, "No programmers there"))
		return false
	}
	if err != nil {
		p.bot.Send(tgbotapi.NewMessage(chat_id, "Something wrong with db, sorry"))
		return false
	}

	matchesArr := strings.Split(matches, " ")

getRespondLoop1:
	for len(matchesArr) != 0 {
		temp, err := strconv.Atoi(matchesArr[0])
		if err != nil {
			matchesArr = matchesArr[1:]
			continue
		}
		user, err := p.db.GetUser(temp)
		if err != nil {
			p.bot.Send(tgbotapi.NewMessage(chat_id, "Can not get user"))
			return false
		}
		p.showProfile(chat_id, user)
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageShowMatchesMenu))
	getRespondLoop2:
		for {
			select {
			case <-timer.C:
				return true
			case <-advertTimer.C:
				if p.advert != nil {
					isAdvertTimer = true
				}
			case upd := <-updates:
				if upd.Message != nil && len(upd.Message.Text) == 1 {
					p.resetTimer(timer)

					switch upd.Message.Text {
					case "1":
						var ChatInfoConf tgbotapi.ChatInfoConfig
						ChatInfoConf.ChatID = int64(user.ChatID)
						chat, err := p.bot.GetChat(ChatInfoConf)
						if err != nil {
							p.bot.Send(tgbotapi.NewMessage(chat_id, "Something wrong with telegram, sorry("))
						}

						if chat.UserName != "" {
							p.bot.Send(tgbotapi.NewMessage(chat_id, "@"+chat.UserName))
						} else {
							var ChatInfoConf tgbotapi.ChatInfoConfig
							ChatInfoConf.ChatID = int64(chat_id)
							chat2, err := p.bot.GetChat(ChatInfoConf)
							if err != nil {
								p.bot.Send(tgbotapi.NewMessage(chat_id, "Something wrong with telegram, sorry("))
							}
							if chat2.UserName != "" {
								p.bot.Send(tgbotapi.NewMessage(chat_id, "The bot will send to this user your username, because this user does not have a username"))
								p.bot.Send(tgbotapi.NewMessage(int64(user.ChatID), "This user has matched you back --> @"+chat2.UserName))
							} else {
								p.bot.Send(tgbotapi.NewMessage(chat_id, "")) // надо эту штуку доделать
							}
						}

						matchesArr = matchesArr[1:]
						break getRespondLoop2
					case "2":
						matchesArr = matchesArr[1:]
						break getRespondLoop2
					case "4":
						if isAdvertTimer {
							p.showAd(chat_id, p.advert, updates, timer)
							p.resetAdvertTimer(advertTimer)
							isAdvertTimer = false
						}
						break getRespondLoop1
					default:
						p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
					}
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
				}
			}
		}
		if isAdvertTimer {
			p.showAd(chat_id, p.advert, updates, timer)
			p.resetAdvertTimer(advertTimer)
			isAdvertTimer = false
		}
	}
	var matchesLeft string
	for _, i := range matchesArr {
		matchesLeft += " " + i
	}
	p.db.SetMatches(chat_id, matchesLeft)
	return false
}
