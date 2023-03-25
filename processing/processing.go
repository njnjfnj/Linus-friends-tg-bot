package processing

import (
	"LinusFriends/LinusUser"
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"fmt"
	"strconv"
	"strings"

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
				switch check {
				case 1:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "SEARCHING..."))
				case 2:
					user, err := p.db.GetUser(int(chat_id))
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "Can not get user"+err.Error()))
					}
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

func (p *Processing) showProfile(user LinusUser.User) {
	profile := tgbotapi.NewPhoto(int64(user.ChatID), tgbotapi.FileBytes{
		Bytes: user.Image,
	})

	profile.Caption = user.Name + "\n" +
		user.SkillsString +
		"\nProgramming experience: " + fmt.Sprintf("%f years\n", user.YearsOfProgramming) +
		user.Description + "\n"

	p.bot.Send(profile)
}

func (p *Processing) showProfileMenu(chat_id int64, updates chan tgbotapi.Update, user LinusUser.User) {
	for true {
		p.showProfile(user)
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
							p.bot.Send(tgbotapi.NewMessage(chat_id, "It is not a number!!"))
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

						user.SkillsMap = make(map[string]bool)
						user.SkillsString = buf
						for _, i := range strings.Split(buf, " ") {
							user.SkillsMap[i] = true
						}
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
