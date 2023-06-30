package processing

import (
	"LinusFriends/storage"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) showMatches(chat_id int64, updates chan tgbotapi.Update, timer *time.Timer) bool {
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
						p.bot.Send(tgbotapi.NewMessage(chat_id, "@"+chat.UserName))
						matchesArr = matchesArr[1:]
						break getRespondLoop2
					case "2":
						matchesArr = matchesArr[1:]
						break getRespondLoop2
					case "4":

						break getRespondLoop1
					default:
						p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
					}
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
				}
			}
		}
	}
	var matchesLeft string
	for _, i := range matchesArr {
		matchesLeft += " " + i
	}
	p.db.SetMatches(chat_id, matchesLeft)
	return false
}
