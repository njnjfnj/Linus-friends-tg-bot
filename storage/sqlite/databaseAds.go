package sqlite

import (
	"LinusFriends/advertisement"
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"database/sql"
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
	q := `INSERT INTO ads (advert_id, content, rating, seen, rated, description) VALUES(?, ?, ?, ?, ?, ?)`

	//ratedratedratedratedratedrated

	if _, err := d.db.ExecContext(d.cntxt, q,
		ad.Advert_id,
		ad.Content,
		ad.Rate,
		ad.Seen,
		ad.Rated,
		ad.Description); err != nil {
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

func (d *Database) UpdateAdContent(content []byte, description string, advert_id int) error {
	q := `UPDATE ads SET content = ?, description = ? WHERE advert_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, content, description, advert_id); err != nil {
		return e.Wrap("can not update content of advert", err)
	}
	return nil
}

func (d *Database) UpdateAdRatingAndViews(rating int, seen int, rated int, advert_id int) error {
	q := `UPDATE ads SET rating = ?, seen = ?, rated = ? WHERE advert_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, rating, seen, rated, advert_id); err != nil {
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
		&res.Seen,
		&res.Rated,
		&res.Description)

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
	rows, err := d.db.Query(q)
	defer func() {
		if err1 := rows.Close(); err1 != nil {
			err = e.Wrap(err1.Error(), err)
		}
	}()
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoAds
	}
	if err != nil {
		return nil, e.Wrap("can not get adverts IDs", err)
	}

	for rows.Next() {
		var buf int
		if err := rows.Scan(&buf); err != nil {
			return nil, e.Wrap("can not scan adverts IDs", err)
		}
		res = append(res, buf)
	}

	return res, nil
}

func (d *Database) IsAdExists(advert_id int) (bool, error) {
	q := `SELECT COUNT(*) FROM ads WHERE advert_id = ?`
	var res int

	if err := d.db.QueryRowContext(d.cntxt, q, advert_id).Scan(&res); err != nil {
		return false, e.Wrap("Can not check if advert exists", err)
	}
	return res > 0, nil
}
