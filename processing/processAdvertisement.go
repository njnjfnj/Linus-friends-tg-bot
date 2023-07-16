package processing

import (
	"LinusFriends/advertisement"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processing) addAdvert(chat_id int64, updates chan tgbotapi.Update) {
	position := 0
	var res advertisement.Ad
	p.bot.Send(tgbotapi.NewMessage(chat_id, "To cancel enter c"))
	p.bot.Send(tgbotapi.NewMessage(chat_id, "Send me a photo of advert"))

forLoop1:
	for upd := range updates {
		if upd.Message != nil {
			switch position {
			case 0:
				if upd.Message.Text == "c" {
					break forLoop1
				}
				if upd.Message.Photo != nil {
					bufContentImage, err := p.processImage(chat_id, upd)
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "error"))
						continue
					}

					res.Content = bufContentImage

					p.bot.Send(tgbotapi.NewMessage(chat_id, "Send me a description of advert"))
					position++

				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "error"))
				}
			case 1:
				if len(upd.Message.Text) != 0 {
					if upd.Message.Text == "c" {
						break forLoop1
					}
					res.Description = upd.Message.Text

					p.bot.Send(tgbotapi.NewMessage(chat_id, "Send me an id (int) of advert"))
					position++
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "error"))
				}
			case 2:
				if len(upd.Message.Text) != 0 {
					if upd.Message.Text == "c" {
						break forLoop1
					}
					bufId, err := strconv.Atoi(upd.Message.Text)
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "It is not a number"))
						continue
					}
					isExists, err := p.db.IsAdExists(bufId)
					if err != nil {
						p.bot.Send(tgbotapi.NewMessage(chat_id, err.Error()))
						continue
					}

					if isExists {
						p.bot.Send(tgbotapi.NewMessage(chat_id, "This id is exists"))
						continue
					}

					res.Advert_id = bufId
					res.Rate = 3
					res.Rated = 1
					res.Seen = 1
					p.processAdvertMessage(chat_id, &res)
					p.bot.Send(tgbotapi.NewMessage(chat_id, "1 - confirm advert\n2 - discard new advert"))
					position++
				} else {
					p.bot.Send(tgbotapi.NewMessage(chat_id, "error"))
				}
			case 3:
				if len(upd.Message.Text) == 1 {
					switch upd.Message.Text {
					case "1":
						if err := p.db.AddNewAd(res); err != nil {
							p.bot.Send(tgbotapi.NewMessage(chat_id, "Error: "+err.Error()))

							p.bot.Send(tgbotapi.NewMessage(chat_id, "1 - confirm advert\n2 - discard new advert"))
							continue
						}

						break forLoop1
					case "2":
						break forLoop1
					default:
						continue
					}
				}
			}
		}
	}
}

func (p *Processing) processAdvertMessage(target_id int64, MessageAdvert *advertisement.Ad) {
	advert := tgbotapi.NewPhoto(target_id, tgbotapi.FileBytes{
		Bytes: MessageAdvert.Content,
	})
	advert.Caption = MessageAdvert.Description

	if _, err := p.bot.Send(advert); err != nil {
		p.bot.Send(tgbotapi.NewMessage(target_id, "Can not show advert: "+err.Error()))
	}
}

func (p *Processing) ShowAds(chat_id int64, updates chan tgbotapi.Update) {
	adverts, err := p.db.GetAdsIds()
	if err != nil {
		p.bot.Send(tgbotapi.NewMessage(chat_id, "can not get ids: "+err.Error()))
		return
	}
	if len(adverts) == 0 {
		p.bot.Send(tgbotapi.NewMessage(chat_id, "There are no adverts("))
		return
	}
	index := 0
	maxIndex := len(adverts) - 1
loop1:
	for {
		advert, err := p.db.GetAd(adverts[index])
		if err != nil {
			p.bot.Send(tgbotapi.NewMessage(chat_id, "can not get advert: "+err.Error()))
			return
		}
		p.processAdvertMessage(chat_id, &advert)
		p.bot.Send(tgbotapi.NewMessage(chat_id, MAAdvertMenu))

	responseLoop1:
		for upd := range updates {
			if len(upd.Message.Text) != 0 {
				switch upd.Message.Text {
				case "1":
					if index == maxIndex {
						index = 0
						break responseLoop1
					}
					index++
					break responseLoop1
				case "2":
					if index == 0 {
						index = maxIndex
						break responseLoop1
					}
					index--
					break responseLoop1
				case "3":
					if p.showChangeAddMenu(chat_id, updates, advert) { // if add has been deleted
						if len(adverts) == 1 {
							return
						} else if index == maxIndex && len(adverts) > 1 {
							index--
						} else {
							index++
						}
						adverts = append(adverts[:index], adverts[:index+1]...)
					}
					break responseLoop1
				case "0":
					break loop1
				case "Post advert":
					p.advert = &advert
					p.bot.Send(tgbotapi.NewMessage(chat_id, "All right!"))
					p.bot.Send(tgbotapi.NewMessage(chat_id, MAAdvertMenu))
				default:
					p.bot.Send(tgbotapi.NewMessage(chat_id, "error"))
					p.bot.Send(tgbotapi.NewMessage(chat_id, MAAdvertMenu))
				}

			}
		}
	}
}

func (p *Processing) showChangeAddMenu(chat_id int64, updates chan tgbotapi.Update, ad advertisement.Ad) bool {
	p.bot.Send(tgbotapi.NewMessage(chat_id, MAAdvertChangeMenu))
Loop1:
	for upd := range updates {
		if len(upd.Message.Text) != 0 {
			switch upd.Message.Text {
			case "0":
				if err := p.db.UpdateAdContent(ad.Content, ad.Description, ad.Advert_id); err != nil {
					p.bot.Send(tgbotapi.NewMessage(chat_id, err.Error()))
				}
				break Loop1
			case "1":
				p.bot.Send(tgbotapi.NewMessage(chat_id, "Send new descriprion (type c to cancel)"))
			loop2:
				for upd := range updates {
					if len(upd.Message.Text) != 0 {
						if upd.Message.Text == "c" {
							break loop2
						}
						ad.Description = upd.Message.Text
						break loop2
					}
				}
				p.bot.Send(tgbotapi.NewMessage(chat_id, MAAdvertChangeMenu))
			case "2":
				p.bot.Send(tgbotapi.NewMessage(chat_id, "Send new photo (type c to cancel)"))
			loop3:
				for upd := range updates {
					if len(upd.Message.Text) != 0 && upd.Message.Text == "c" {
						break loop3
					}
					if upd.Message.Photo != nil {
						buf, err := p.processImage(chat_id, upd)
						if err != nil {
							p.bot.Send(tgbotapi.NewMessage(chat_id, MessageErrorCanNotUploadPhoto))
							continue loop3
						}
						ad.Content = buf
						break loop3
					}
				}
				p.bot.Send(tgbotapi.NewMessage(chat_id, MAAdvertChangeMenu))
			case "Delete advertisement":
				p.db.DeleteAd(ad.Advert_id)
				return true
			default:
				p.bot.Send(tgbotapi.NewMessage(chat_id, "error"))
				p.bot.Send(tgbotapi.NewMessage(chat_id, MAAdvertChangeMenu))
			}
		}
	}
	return false
}
