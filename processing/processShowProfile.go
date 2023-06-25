package processing

import (
	"LinusFriends/LinusUser"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) showProfile(target_id int64, user LinusUser.User) {
	profile := tgbotapi.NewPhoto(target_id, tgbotapi.FileBytes{
		Bytes: user.Image,
	})

	profile.Caption = user.Name + "\n" +
		user.SkillsString +
		"\nProgramming experience: " + fmt.Sprintf("%f years\n", user.YearsOfProgramming) +
		user.Description + "\n"

	_, err := p.bot.Send(profile)
	if err != nil {
		p.bot.Send(tgbotapi.NewMessage(target_id, err.Error()))
	}
}

func (p *Processing) showProfileMenu(chat_id int64, updates chan tgbotapi.Update, user LinusUser.User) {
	for {
		p.showProfile(int64(user.ChatID), user)
		p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeProfile))
		var check int8
		check = -1
		for upd := range updates {
			if upd.Message != nil {
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
						p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeSkills))
					case 6:
						p.db.UpdateUser(user)
						return
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
					case 5:
						if len(upd.Message.Text) > 200 {
							p.bot.Send(tgbotapi.NewMessage(chat_id, "Maximum length - 200 characters"))
							continue
						}

						buf := strings.ToLower(upd.Message.Text)

						user.SkillsString = buf
						check = -1
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
			}
			if check == -1 {
				p.bot.Send(tgbotapi.NewMessage(chat_id, "Successfull local change!!!"))
				break
			}
		}
	}
}
