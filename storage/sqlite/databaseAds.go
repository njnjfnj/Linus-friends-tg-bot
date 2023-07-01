package sqlite

import (
	"LinusFriends/advertisement"
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"database/sql"
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
короче, инитишь дб также, как и дб рользователей,
но дб пользователей вроде как лучше маршалить, но это не точно.

Сделаешь отдельную горутину для рассылки рекламы,
а в режиме админа сделаешь так, что бы можно было добавлять рекламу
и смотреть ее, рейтинг и сколько ее просмотрело.

Механизм рассылки рекламы будет примерно таким: через режим админа
запуск рассылки, после этого извлекается массив id, потом по очереди,
с интервалом во времени, который можно будет настроить, будет вестись рассылка.
*/

func (d *Database) AddNewAd(ad advertisement.Ad) error {
	q := `INSERT INTO ads (advert_id, content, rating, seen) VALUES(?, ?, ?, ?)`

	bufContent, err := json.Marshal(ad.Content)
	if err != nil {
		return e.Wrap("can not parse to json", err)
	}
	//ratedratedratedratedratedrated

	if _, err := d.db.ExecContext(d.cntxt, q,
		ad.Advert_id,
		bufContent,
		ad.Rate,
		ad.Seen); err != nil {
		return e.Wrap("Can not add new advertisement", err)
	}

	return nil
}

func (d *Database) DeleteAd(advert_id int) error {
	q := `DELETE FROM ads WHERE advert_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, advert_id); err != nil {
		return e.Wrap("can not delete advert", err)
	}
	return nil
}

func (d *Database) UpdateAdContent(content *tgbotapi.PhotoConfig, advert_id int) error {
	q := `UPDATE ads SET content = ? WHERE advert_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, *content, advert_id); err != nil {
		return e.Wrap("can not update content of advert", err)
	}
	return nil
}

func (d *Database) UpdateAdRatingAndViews(rating float32, seen int, advert_id int) error {
	q := `UPDATE ads SET rating = ?, seen = ? WHERE advert_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, rating, seen, advert_id); err != nil {
		return e.Wrap("can not update rating and views of advert", err)
	}
	return nil
}

func (d *Database) GetAd(advert_id int) (res advertisement.Ad, err error) {
	q := `SELECT * FROM ads WHERE advert_id = ?`

	err = d.db.QueryRowContext(d.cntxt, q, advert_id).Scan(
		&res.Advert_id,
		&res.Content,
		&res.Rate,
		&res.Seen)

	if err == sql.ErrNoRows {
		return advertisement.Ad{}, storage.ErrNoAds
	}
	if err != nil {
		return advertisement.Ad{}, e.Wrap("can not get advert", err)
	}

	return res, nil
}

func (d *Database) GetAds() (res []advertisement.Ad, err error) {
	q := `SELECT * FROM ads`
	err = d.db.QueryRowContext(d.cntxt, q).Scan(&res)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoAds
	}
	if err != nil {
		return nil, e.Wrap("can not get adverts", err)
	}

	return res, nil
}

func (d *Database) GetAdsIds() (res []int, err error) {
	q := `SELECT advert_id FROM ads`
	err = d.db.QueryRowContext(d.cntxt, q).Scan(&res)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoAds
	}
	if err != nil {
		return nil, e.Wrap("can not get adverts IDs", err)
	}

	return res, nil
}
