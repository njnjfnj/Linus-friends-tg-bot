package advertisement

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Ad struct {
	Advert_id string
	Content   tgbotapi.PhotoConfig
	seen      int
	Rate      uint8
}
