package processing

import (
	"LinusFriends/LinusUser"
	"LinusFriends/storage"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) searchForProgrammers(chat_id int64, updates chan tgbotapi.Update, user LinusUser.User) {
	for {
		countOfErrors, countOfHits := 0, 0
		var searchingByWhat int
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageSearchForProgrammersMenu))
	getRespondLoop1:
		for upd := range updates {
			if upd.Message != nil && len(upd.Message.Text) == 1 {
				check, err := strconv.Atoi(upd.Message.Text)
				if err != nil {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "It is not a number!!"))
					continue
				}
				switch check {
				case 1:
					searchingByWhat = storage.SearchingByExperience
					break getRespondLoop1
				case 2:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "This option is currently under development"))
					continue
					//searchingByWhat = storage.SearchingByLanguage
				case 3:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "This option is currently under development"))
					continue
					//searchingByWhat = storage.SearchingByLanguagesAndExpirience
				case 4:
					searchingByWhat = storage.SearchingByRandom
					break getRespondLoop1
				case 5:
					return
				default:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 4]!!!! 🤬🤬"))
				}
			} else {
				p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 4]!!!! 🤬🤬"))
			}
		}
	getRespondLoop2:
		for true {
			friend, err := p.db.GetRandomUserForUser(chat_id, searchingByWhat, user)
			if err == storage.ErrNoFriends {
				p.bot.Send(tgbotapi.NewMessage(chat_id, "No friends have found\n🥲😅😂🤣🤣😓😪😭😭😭😭😭\nTry another searching method or try again later"))
				continue
			} else if err != nil {
				countOfErrors++
				p.bot.Send(tgbotapi.NewMessage(chat_id, "ERROR: can not get user\n"+err.Error()))
				if countOfErrors > 6 {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "Please, try again later"))
					return
				}
				continue
			} else if countOfErrors != 0 {
				countOfHits++
				if countOfHits > 2 {
					countOfErrors, countOfHits = 0, 0
				}
			}
			p.showProfile(friend)
			p.bot.Send(tgbotapi.NewMessage(chat_id, MessageIntaractionWithFriend))
		getRespondLoop3:
			for upd := range updates {
				if upd.Message != nil && len(upd.Message.Text) == 1 {
					check, err := strconv.Atoi(upd.Message.Text)
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "It is not a number!!"))
						continue
					}
					switch check {
					case 1:
						break getRespondLoop3
					case 2:
						break getRespondLoop3
					case 4:
						break getRespondLoop2
					default:
						p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
					}
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
				}
			}
		}
		searchingByWhat = -1
	}
}
