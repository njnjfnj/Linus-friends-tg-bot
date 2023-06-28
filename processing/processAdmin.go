package processing

import (
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) processAdmin(chat_id int64, updates chan tgbotapi.Update, timer *time.Timer) {
	p.bot.Send(tgbotapi.NewMessage(chat_id, "私̸ͪ͠_͕̚͡はあ̶̵̴̧̛̜̬̼̳̭̱̦̬̜̰̫̲̮̋̑ͫ́̓̎̔ͮ͒ͮ̓̌͟͝な̸̵̭̼̘͉̿̐͂ͯͅ_̴̧̡̢̧̟̫̹̰̜̬͓ͥ̅̽ͥ̑̅͆̋̋͘͡͝_̴ͮた̶̙̾͌̋̈́͟を̷̷̧̛̪͈̬̙͎̻͈͇̜̩ͮ̓ͥ̔ͣͪ͒̊́ͫ̓ͩͧ͗̍͐̈̔̂͜ͅ見̶̖̩̦̗̪̗̟̺̣̥̥̗̤͕͕͎͗͒ͭ̅ͪͫ͆̅ͩ̈̈ͨ̓ͧͭ͐̂̓ͬ́̈́̀̀͜͝͝͠͠ͅつ̤͕̉ͯ̍け̵̸̙͇̭͇̱̣̜ͤ̓̿ͦ͒ͮ̾ͥ͢_̢̧͙̯̬̰̖͔̔͒ͦ͑̅̚る̴̡̛̛̙̯͖̺̗͍̝͓͈̖̩̦̖͈̰͊͗ͭ͆ͣ̃̐͋͑ͬ͋͑́̑ͮ͘̚̕̕͝で̵̶̛̤̬̻̣͈̫̺͚͇̞͌ͬ̓̿͂ͪ̋̓ͦ̽ͦͨ̏͑̐͊ͬ͌͗͘̚͜し̧̰̣͚͓͓͇̺͔̞́̑ͣ̎ͤ̇̎ͫ͗ͣ̕ͅょ̴̪̙͎̣̞̳͙͕́͂̐ͬ̏͒͒̂̐ͯ̃̌ͮ̈́̒ͫ͂̕͡う̗̩̺͚"))

getResponseLoop1:
	for {
		select {
		case <-timer.C:
			return
		case upd := <-updates:
			if upd.Message != nil {
				if upd.Message.Text != "" && upd.Message.Text == p.adminPassword {
					break getResponseLoop1
				} else {
					return
				}
			}
		}
	}

getResponseLoop2:
	for {
		p.bot.Send(tgbotapi.NewMessage(chat_id, MAMenu))
		for upd := range updates {
			if upd.Message != nil {
				switch upd.Message.Text { // от сделай везде, что бы в инт не парсилось, а то херню сделал
				case "0":
					break getResponseLoop2
				case "1":

				case "2":

				case "3":
					number, err := p.db.UserCount()
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, err.Error()))
					}
					p.bot.Send(tgbotapi.NewMessage(chat_id, strconv.Itoa(number)))
				default:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "incorrect"))
				}
				p.bot.Send(tgbotapi.NewMessage(chat_id, MAMenu))
			}
		}
	}

	p.resetTimer(timer)
}
