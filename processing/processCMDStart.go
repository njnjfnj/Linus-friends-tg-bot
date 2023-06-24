package processing

import (
	"LinusFriends/LinusUser"
	"LinusFriends/libs/e"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) processCMDStart(chat_id int64, updates chan tgbotapi.Update) error {
	check, err := p.db.IsUserExists(int(chat_id))
	if err != nil {
		return e.Wrap("Can not chek if user exists", err)
	}
	if check {
		p.showMenu(chat_id, updates)
		return nil
	}
	p.bot.Send(tgbotapi.NewMessage(chat_id, MessageStart+"\n\n"+MessageChangeName))

	var newUser LinusUser.User
	newUser.ChatID = int(chat_id)
	position := 0

	for upd := range updates {
		if upd.Message != nil { // upd.FromChat().ChatConfig().ChatID == chat_id &&
			switch position {
			case 0:
				if len(upd.Message.Text) != 0 {
					if len(upd.Message.Text) > 50 {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "Maximum name length - 50 characters"))
						continue
					}
					newUser.Name = upd.Message.Text
					p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeDescription))
					position++
				}
			case 1:
				if len(upd.Message.Text) != 0 {
					if len(upd.Message.Text) > 1500 {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "Maximum description length - 1500 characters"))
						continue
					}
					newUser.Description = upd.Message.Text
					p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeProfilePic))
					position++
				}
			case 2:
				if upd.Message.Photo != nil {
					if newUser.Image, err = p.processImage(chat_id, upd); err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, MessageErrorCanNotUploadPhoto))
						log.Println("Can not process a cmd: ", err)
						continue
					}

					p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeSkills))
					position++
				}
			case 3:
				if len(upd.Message.Text) != 0 {
					if len(upd.Message.Text) > 200 {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "Maximum length - 200 characters"))
						continue
					}

					buf := strings.ToLower(upd.Message.Text)

					//newUser.SkillsMap = make(map[string]bool)
					newUser.SkillsString = buf
					// for _, i := range strings.Split(buf, " ") {
					// 	newUser.SkillsMap[i] = true
					// }

					p.bot.Send(tgbotapi.NewMessage(chat_id, MessageChangeYearsPfProgramming))
					position++
				}
			case 4:
				if len(upd.Message.Text) != 0 {
					buf, err := strconv.ParseInt(upd.Message.Text, 0, 64)
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "Can not parse to int"))
						continue
					}
					newUser.YearsOfProgramming = int(buf)
					position++
					//break
				}
			}
			if position == 5 {
				break
			}
		}
	}
	if err := p.db.AddNewUser(newUser); err != nil {
		p.bot.Send(tgbotapi.NewMessage(chat_id, "Can not add new user. Try again later (use /start command): "+err.Error()))
		return err
	}

	p.showMenu(chat_id, updates)

	return nil
}

func (p *Processing) processImage(chat_id int64, upd tgbotapi.Update) (buf []byte, err error) {
	defer func() { err = e.WrapIfErr("Can not get image", err) }()
	f, err := p.bot.GetFile(tgbotapi.FileConfig{
		FileID: upd.Message.Photo[len(upd.Message.Photo)-1].FileID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		f.Link(p.bot.Token),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := p.bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
