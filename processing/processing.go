package processing

import (
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"LinusFriends/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processing struct {
	bot *tgbotapi.BotAPI
	db  storage.Storage
}

func (p *Processing) CMD(cmd string, chat_id int64) (err error) {

	// ЕБААААТЬ, короче можно через get update chan замутить бизнесяку ахуенную и не ебатся с хуерлышами поганками

	defer func() { err = e.WrapIfErr("Can not process a cmd: ", err) }()
	switch cmd {
	case "/help":
		_, err = p.bot.Send(tgbotapi.NewMessage(chat_id, MessageHelp))
		if err != nil {
			return err
		}
	case "/start":
		_, err = p.bot.Send(tgbotapi.NewMessage(chat_id, MessageStart+"\n"+MessageChangeName))
		if err != nil {
			return err
		}

		if err := p.db.AddNewUser(user.User{
			ChatID:      int(chat_id),
			LastCommand: user.CmdChangeProfileName,
			IsImportant: true,
		}); err != nil {
			return err
		}
	default:
		_, err = p.bot.Send(tgbotapi.NewMessage(chat_id, "Unknown command"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processing) Photo(photo []tgbotapi.PhotoSize, chat_id int64) error {
	// usr, err := p.db.GetUser(int(chat_id))

	// if err != nil {

	// }

	// if usr.LastCommand == user.CmdChangeProfile || usr.LastCommand == user.CmdChangeProfilePhoto {

	// }
	return nil
}

func (p *Processing) Text(text string, chat_id int64) error {

	return nil
}

func NewProcessing(botapi *tgbotapi.BotAPI, storage storage.Storage) *Processing {
	return &Processing{
		bot: botapi,
		db:  storage,
	}
}
