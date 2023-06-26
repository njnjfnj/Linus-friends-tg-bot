package sqlite

import (
	"LinusFriends/LinusUser"
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"context"
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db    *sql.DB
	cntxt context.Context
}

func NewDatabase(path string) (*Database, error) {
	DB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can not open DB", err)
	}
	if err := DB.Ping(); err != nil {
		return nil, e.Wrap("can not connect to DB", err)
	}
	return &Database{db: DB}, nil
}

func (d *Database) Init(context context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS db (chat_id INT NOT NULL, name TEXT, description TEXT, skillsString TEXT NOT NULL, years_of_programming INT NOT NULL, photo BLOB NOT NULL); CREATE TABLE IF NOT EXISTS skls (language TEXT NOT NULL, IDs TEXT); CREATE TABLE IF NOT EXISTS mtchs (chat_id INT NOT NULL, IDs TEXT NOT NULL); CREATE UNIQUE INDEX IF NOT EXISTS idx_id ON db (chat_id); CREATE UNIQUE INDEX IF NOT EXISTS idx_lng ON skls (language); CREATE UNIQUE INDEX IF NOT EXISTS idx_id ON mtchs (chat_id);` // (STRING, BIT))
	if _, err := d.db.ExecContext(context, q); err != nil {
		return e.Wrap("Can not create DB", err)
	}
	d.cntxt = context

	return nil
}

func (d *Database) AddNewUser(u LinusUser.User) (err error) {
	defer func() { err = e.WrapIfErr("Can not add new user", err) }()

	q5 := `INSERT INTO mtchs (chat_id, IDs) VALUES(?, ?)`
	if _, err := d.db.ExecContext(d.cntxt, q5, u.ChatID, ""); err != nil {
		return e.Wrap("Can not add to matches db", err)
	}

	q2 := `SELECT COUNT(*) FROM skls WHERE language = ?`
	q3 := `UPDATE skls SET IDs = IDs || ' ' || ? WHERE language = ?`
	q4 := `INSERT INTO skls (language, IDs) VALUES(?, ?)`
	var res int
	for _, i := range strings.Split(u.SkillsString, " ") {
		if err := d.db.QueryRowContext(d.cntxt, q2, i).Scan(&res); err != nil {
			return e.Wrap("Can not check if language exists", err)
		} else if res > 0 {
			if _, err := d.db.ExecContext(d.cntxt, q3, strconv.Itoa(u.ChatID), i); err != nil {
				return e.Wrap("Can not update skills IDs", err)
			}
		} else if _, err := d.db.ExecContext(d.cntxt, q4, i, strconv.Itoa(u.ChatID)); err != nil {
			return e.Wrap("Can not add new language", err)
		}

	}

	q1 := `INSERT INTO db (chat_id, name, description, skillsString, years_of_programming, photo) VALUES(?, ?, ?, ?, ?, ?)`
	if _, err := d.db.ExecContext(d.cntxt, q1,
		u.ChatID,
		u.Name,
		u.Description,
		u.SkillsString,
		u.YearsOfProgramming,
		u.Image); err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteUser(chat_id int) error {
	q := `DELETE FROM db WHERE chat_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, chat_id); err != nil {
		return e.Wrap("Can not delete user", err)
	}
	return nil
}

func (d *Database) IsUserExists(chat_id int) (bool, error) {
	q := `SELECT COUNT(*) FROM db WHERE chat_id = ?`
	var res int

	if err := d.db.QueryRowContext(d.cntxt, q, chat_id).Scan(&res); err != nil {
		return false, e.Wrap("Can not check if user exists", err)
	}
	return res > 0, nil
}

func (d *Database) GetUser(chat_id int) (LinusUser.User, error) {
	var res LinusUser.User
	q := `SELECT * FROM db WHERE chat_id = ?`
	err := d.db.QueryRowContext(d.cntxt, q, chat_id).Scan(
		&res.ChatID,
		&res.Name,
		&res.Description,
		&res.SkillsString,
		&res.YearsOfProgramming,
		&res.Image)

	if err == sql.ErrNoRows {
		return LinusUser.User{}, storage.ErrNoFriends
	}
	if err != nil {
		return LinusUser.User{}, e.Wrap("Can not get user", err)
	}

	return res, nil
}

func (d *Database) UpdateUser(u LinusUser.User) error {
	q := `UPDATE db SET name = ?, description = ?, skillsString = ?, years_of_programming = ?, photo = ? WHERE chat_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q,
		u.Name,
		u.Description,
		u.SkillsString,
		u.YearsOfProgramming,
		u.Image,
		u.ChatID); err != nil {
		return e.Wrap("Can not update user", err)
	}
	return nil
}

func (d *Database) GetRandomUserForUser(chat_id int64, SearchByWhat int, user LinusUser.User) (LinusUser.User, string, error) {
	var res LinusUser.User
	var ids = ""
	var err error
	switch SearchByWhat {
	case storage.SearchingByExperience:
		q := `SELECT * FROM db WHERE years_of_programming = ? EXCEPT SELECT * FROM db WHERE chat_id = ?`
		err = d.db.QueryRowContext(d.cntxt, q, user.YearsOfProgramming, chat_id).Scan(
			&res.ChatID,
			&res.Name,
			&res.Description,
			&res.SkillsString,
			&res.YearsOfProgramming,
			&res.Image)
	case storage.SearchingByRandom:
		q := `SELECT * FROM (SELECT * FROM db  EXCEPT SELECT * FROM db WHERE chat_id = ?) ORDER BY RANDOM() LIMIT 1`
		err = d.db.QueryRowContext(d.cntxt, q, chat_id).Scan(
			&res.ChatID,
			&res.Name,
			&res.Description,
			&res.SkillsString,
			&res.YearsOfProgramming,
			&res.Image)
	case storage.SearchingByLanguage:
		q := `SELECT IDs FROM skls WHERE language = ?`
		var temp string
		for _, i := range strings.Split(user.SkillsString, " ") {
			err = d.db.QueryRowContext(d.cntxt, q, i).Scan(&temp)
			if err != nil {
				break
			}
			ids += temp + " "
		}
	}

	if err == sql.ErrNoRows {
		return LinusUser.User{}, "", storage.ErrNoFriends
	}
	if err != nil {
		return LinusUser.User{}, "", e.Wrap("Can not get random user", err)
	}

	return res, ids, nil
}

func (d *Database) GetMatches(chat_id int64) (res string, err error) {
	defer func() { err = e.WrapIfErr("Can not get matches", err) }()

	q := `SELECT IDs FROM mtchs WHERE chat_id = ?`
	err = d.db.QueryRowContext(d.cntxt, q, chat_id).Scan(&res)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoFriends
	}
	if err != nil {
		return "", e.Wrap("Can not do query", err)
	}

	return res, err
}

func (d *Database) SetMatches(chat_id int64, matchesLeft string) error {
	q := `UPDATE mtchs SET chat_id = ?, IDs = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, chat_id, matchesLeft); err != nil {
		return e.Wrap("Can not set matches", err)
	}
	return nil
}

func (d *Database) AddMatch(chat_id int64, u LinusUser.User) error {
	q := `UPDATE mtchs SET IDs = IDs || ' ' || ? WHERE chat_id = ?`
	if _, err := d.db.ExecContext(d.cntxt, q, strconv.Itoa(u.ChatID), chat_id); err != nil {
		return e.Wrap("Can not add user to matches", err)
	}
	return nil
}
