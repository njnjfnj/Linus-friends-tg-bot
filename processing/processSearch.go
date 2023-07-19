package processing

import (
	"LinusFriends/LinusUser"
	"LinusFriends/storage"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) searchForProgrammers(chat_id int64, updates chan tgbotapi.Update, user LinusUser.User, timer *time.Timer, advertTimer *time.Timer) bool {
	isAdvertTimer := false
	for {
		countOfErrors, countOfHits := 0, 0
		var searchingByWhat int
	getRespondLoop1:
		for {
			p.bot.Send(tgbotapi.NewMessage(chat_id, MessageSearchForProgrammersMenu))
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
						searchingByWhat = storage.SearchingByExperience
						break getRespondLoop1
					case "2":
						searchingByWhat = storage.SearchingByLanguage
						break getRespondLoop1
					case "3":
						searchingByWhat = storage.SearchingByRandom
						break getRespondLoop1
					case "4":
						return false
					default:
						p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 4]!!!! 🤬🤬"))
					}
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 4]!!!! 🤬🤬"))
				}
				if isAdvertTimer {
					p.showAd(chat_id, p.advert, updates, timer)
					p.resetAdvertTimer(advertTimer)
					isAdvertTimer = false
				}
			}
		}
	getRespondLoop2:
		for {
			friend, ids, err := p.db.GetRandomUserForUser(chat_id, searchingByWhat, user)
			if err == storage.ErrNoFriends {
				p.bot.Send(tgbotapi.NewMessage(chat_id, "No friends have found\n🥲😅😂🤣🤣😓😪😭😭😭😭😭\nTry another searching method or try again later"))
				break
			} else if err != nil {
				countOfErrors++
				p.bot.Send(tgbotapi.NewMessage(chat_id, "ERROR: can not get user\n"+err.Error()))
				if countOfErrors > 6 {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "Please, try again later"))
					return false
				}
				continue
			} else if countOfErrors != 0 {
				countOfHits++
				if countOfHits > 2 {
					countOfErrors, countOfHits = 0, 0
				}
			}
			if ids == "" {
				p.showProfile(chat_id, friend)
				p.bot.Send(tgbotapi.NewMessage(chat_id, MessageIntaractionWithFriend))
			getRespondLoop3:
				for {
					select {
					case <-advertTimer.C:
						if p.advert != nil {
							isAdvertTimer = true
						}
					case <-timer.C:
						return true
					case upd := <-updates:
						if upd.Message != nil && len(upd.Message.Text) == 1 {
							p.resetTimer(timer)

							switch upd.Message.Text {
							case "1":
								if err := p.db.AddMatch(int64(friend.ChatID), user); err != nil {
									p.bot.Send(tgbotapi.NewMessage(chat_id, "Sorry, it is something wrong with bot("))
								}
								p.bot.Send(tgbotapi.NewMessage(chat_id, "successfully matched"))
								break getRespondLoop3
							case "2":
								break getRespondLoop3
							case "4":
								break getRespondLoop2
							default:
								p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
							}
						} else {
							p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
						}
						if isAdvertTimer {
							p.showAd(chat_id, p.advert, updates, timer)
							p.resetAdvertTimer(advertTimer)
							isAdvertTimer = false
						}
					}
				}
			} else {
				p.bot.Send(tgbotapi.NewMessage(chat_id, "A new package of users that know the same programming languages as you has been taken"))
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				idsArr := strings.Split(ids, " ")
				l := len(idsArr)
				idsArr = idsArr[:l-1]
				l--
				seenFriends := make(map[int]bool)
				for l != 0 {
					i := r.Intn(l)
					id, _ := strconv.Atoi(idsArr[i])
					if id == user.ChatID || seenFriends[id] {
						idsArr = append(idsArr[:i], idsArr[i+1:]...)
						l--
						continue
					}
					friend, err = p.db.GetUser(id)
					if err != nil {
						countOfErrors++
						p.bot.Send(tgbotapi.NewMessage(chat_id, "ERROR: can not get user\n"+err.Error()))
						if countOfErrors > 6 {
							p.bot.Send(tgbotapi.NewMessage(chat_id, "Please, try again later"))
							return false
						}
						continue
					} else if countOfErrors != 0 {
						countOfHits++
						if countOfHits > 2 {
							countOfErrors, countOfHits = 0, 0
						}
					}
					p.showProfile(chat_id, friend)
					p.bot.Send(tgbotapi.NewMessage(chat_id, MessageIntaractionWithFriend))
				getRespondLoop4:
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
									if err := p.db.AddMatch(int64(friend.ChatID), user); err != nil {
										p.bot.Send(tgbotapi.NewMessage(chat_id, "Sorry, it is something wrong with bot("))
									}
									p.bot.Send(tgbotapi.NewMessage(chat_id, "successfully matched"))
									break getRespondLoop4
								case "2":
									break getRespondLoop4
								case "4":
									break getRespondLoop2
								default:
									p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
								}
							} else {
								p.bot.Send(tgbotapi.NewMessage(chat_id, "∈[1, 2]U[4, 4]!!!! 🤬🤬"))
							}
							if isAdvertTimer {
								p.showAd(chat_id, p.advert, updates, timer)
								p.resetAdvertTimer(advertTimer)
								isAdvertTimer = false
							}
						}
					}

					idsArr = append(idsArr[:i], idsArr[i+1:]...)
					l--
					seenFriends[id] = true
				}
				p.bot.Send(tgbotapi.NewMessage(chat_id, "There are no friends anymore, you can start this searching method again to find more friends"))
				break getRespondLoop2
			}
		}
		searchingByWhat = -1
	}
}
