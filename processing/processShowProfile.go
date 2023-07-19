package processing

import (
	"LinusFriends/LinusUser"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) showProfile(target_id int64, user LinusUser.User) {
	profile := tgbotapi.NewPhoto(target_id, tgbotapi.FileBytes{
		Bytes: user.Image,
	})

	profile.Caption = user.Name + "\n" +
		user.SkillsString +
		"\nProgramming experience: " + strconv.Itoa(user.YearsOfProgramming) + " years\n" +
		user.Description + "\n"

	_, err := p.bot.Send(profile)
	if err != nil {
		p.bot.Send(tgbotapi.NewMessage(target_id, err.Error()))
	}
}

func (p *Processing) showProfileMenu(chat_id int64, updates chan tgbotapi.Update, user LinusUser.User, timer *time.Timer, advertTimer *time.Timer) bool {
	for { // there are plenty of shitcode
		isAdvertTimer := false
		p.showProfile(int64(user.ChatID), user)
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeProfile))
		var check int8
		check = -1
	getRespondLoop1:
		for {
			select {
			case <-timer.C:
				return true
			case <-advertTimer.C:
				if p.advert != nil {
					isAdvertTimer = true
				}
			case upd := <-updates:
				if upd.Message != nil {
					p.resetTimer(timer)
					if check == -1 && len(upd.Message.Text) == 1 {
						buf, err := strconv.Atoi(upd.Message.Text)
						if err != nil {
							p.bot.Send(tgbotapi.NewMessage(chat_id, "It is not a number!!"))
							continue
						}
						check = int8(buf)
						switch check {
						case 1:
							p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeProfilePic))
						case 2:
							p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeName))
						case 3:
							p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeYearsPfProgramming))
						case 4:
							p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeDescription))
						case 5:
							if p.showChangeSkillsMenu(chat_id, updates, &user, timer, advertTimer) {
								return true
							}
						case 0:
							p.db.UpdateUser(user)
							return false
						default:
							p.bot.Send(tgbotapi.NewMessage(chat_id, "[1, 6]!!!! 🤬🤬"))
						}
					} else if len(upd.Message.Text) != 0 {
						switch check {
						case 2:
							if len(upd.Message.Text) > 50 {
								p.bot.Send(tgbotapi.NewMessage(chat_id, "Maximum name length - 50 characters"))
								continue
							}
							user.Name = upd.Message.Text
							check = -1
						case 3:
							buf, err := strconv.Atoi(upd.Message.Text)
							if err != nil {
								p.bot.Send(tgbotapi.NewMessage(chat_id, "It is not a number!! 🤬🤬"))
								continue
							}
							user.YearsOfProgramming = buf
							check = -1
						case 4:
							if len(upd.Message.Text) > 1500 {
								p.bot.Send(tgbotapi.NewMessage(chat_id, "Maximum description length - 1500 characters"))
								continue
							}
							user.Description = upd.Message.Text
							check = -1
							// case 5:
							// 	if len(upd.Message.Text) > 200 {
							// 		p.bot.Send(tgbotapi.NewMessage(chat_id, "Maximum length - 200 characters"))
							// 		continue
							// 	}

							// 	// тут не изменяются скилы в бд

							// 	// buf := strings.ToLower(upd.Message.Text)

							// 	// user.SkillsString = buf

							// 	check = -1
						}
					} else if upd.Message.Photo != nil && check == 1 {
						buf, err := p.processImage(chat_id, upd)
						if err != nil {
							p.bot.Send(tgbotapi.NewMessage(chat_id, MessageErrorCanNotUploadPhoto))
							continue
						}
						user.Image = buf
						check = -1
					} else {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "ERROR"))
					}
					if isAdvertTimer {
						p.showAd(chat_id, p.advert, updates, timer)
						p.resetAdvertTimer(advertTimer)
						isAdvertTimer = false
					}
				}
				if check == -1 {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "Successfull local change!!!"))
					break getRespondLoop1
				}
			}
		}
	}
}

func (p *Processing) showChangeSkillsMenu(chat_id int64, updates chan tgbotapi.Update, user *LinusUser.User, timer *time.Timer, advertTimer *time.Timer) bool {
	isAdvertTimer := false
	skills := strings.Split(user.SkillsString, " ")
	index := 0
	maxIndex := len(skills) - 1
responseLoop1:
	for {
		if skills[index] == " " {
			skills = append(skills[:index], skills[:index+1]...)
			continue
		}
		p.bot.Send(tgbotapi.NewMessage(chat_id, skills[index]))
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeSkillsMenu))
		select {
		case <-timer.C:
			return true
		case <-advertTimer.C:
			if p.advert != nil {
				isAdvertTimer = true
			}
		case upd := <-updates:
			if upd.Message != nil && len(upd.Message.Text) != 0 {
				p.resetTimer(timer)
				switch upd.Message.Text {
				case "0":
					break responseLoop1
				case "1":
					if index == maxIndex {
						index = 0
						break
					}
					index++
				case "2":
					if index == 0 {
						index = maxIndex
						break
					}
					index--
				case "3":
					if maxIndex == 0 {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "You must have at least 1 skill"))
						break
					}

					if err := p.db.DeleteSkill(int(chat_id), skills[index]); err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, err.Error()))
						break
					}

					user.SkillsString = strings.ReplaceAll(user.SkillsString, skills[index], "")
					skills = append(skills[:index], skills[:index+1]...)
					p.db.UpdateUser(*user)
				case "4":

				default:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "enter number from 0 to 4"))
				}
				if isAdvertTimer {
					p.showAd(chat_id, p.advert, updates, timer)
					p.resetAdvertTimer(advertTimer)
					isAdvertTimer = false
				}
			}
		}
	}
	return false
}
