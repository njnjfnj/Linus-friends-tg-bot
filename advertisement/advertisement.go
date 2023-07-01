package advertisement

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Ad struct {
	Advert_id int
	Content   tgbotapi.PhotoConfig
	Seen      int
	Rated     int
	Rate      float32
}
