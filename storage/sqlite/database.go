package sqlite

import (
	"LinusFriends/LinusUser"
	"LinusFriends/libs/e"
	"LinusFriends/storage"
	"context"
	"database/sql"
	"encoding/json"

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
	q := `CREATE TABLE IF NOT EXISTS db (chat_id INT, name TEXT, description TEXT, skillsString TEXT, skillsJson TEXT, years_of_programming INT, photo BLOB NOT NULL)`
	if _, err := d.db.ExecContext(context, q); err != nil {
		return e.Wrap("Can not create DB", err)
	}
	d.cntxt = context

	return nil
}

func (d *Database) AddNewUser(u LinusUser.User) (err error) {
	defer func() { err = e.WrapIfErr("Can not add new user", err) }()
	q := `INSERT INTO db (chat_id, name, description, skillsString, skillsJson, years_of_programming, photo) VALUES(?, ?, ?, ?, ?, ?, ?)`
	skillsMapBuf, err := json.Marshal(u.SkillsMap)
	if err != nil {
		return err
	}
	if _, err := d.db.ExecContext(d.cntxt, q,
		u.ChatID,
		u.Name,
		u.Description,
		u.SkillsString,
		string(skillsMapBuf),
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
	err := d.db.QueryRowContext(d.cntxt, q, chat_id).Scan(&res)
	if err != nil {
		return false, e.Wrap("Can not check if user exists", err)
	}
	return res > 0, nil
}

func (d *Database) GetUser(chat_id int) (LinusUser.User, error) {
	var res LinusUser.User
	var bufSkillsMap string

	q := `SELECT * FROM db WHERE chat_id = ?`
	err := d.db.QueryRowContext(d.cntxt, q, chat_id).Scan(
		&res.ChatID,
		&res.Name,
		&res.Description,
		&res.SkillsString,
		&bufSkillsMap,
		&res.YearsOfProgramming,
		&res.Image)

	if err == sql.ErrNoRows {
		return LinusUser.User{}, storage.ErrNoSavedPages
	}
	if err != nil {
		return LinusUser.User{}, e.Wrap("Can not get user", err)
	}
	if err = json.Unmarshal([]byte(bufSkillsMap), &res.SkillsMap); err != nil {
		return LinusUser.User{}, e.Wrap("Can now get user", err)
	}
	return res, nil
}

func (d *Database) UpdateUser(u LinusUser.User) error {
	q := `UPDATE db SET chat_id = ?, name = ?, description = ?, skillsString = ?, skillsJson = ?, years_of_programming = ?, photo = ?`
	skillsMapBuf, err := json.Marshal(u.SkillsMap)
	if err != nil {
		return err
	}
	if _, err := d.db.ExecContext(d.cntxt, q,
		u.ChatID,
		u.Name,
		u.Description,
		u.SkillsString,
		string(skillsMapBuf),
		u.YearsOfProgramming,
		u.Image); err != nil {
		return e.Wrap("Can not update user", err)
	}
	return nil
}

func (d *Database) GetRandomUser() (LinusUser.User, error) {
	var res LinusUser.User

	// q := `SELECT * FROM db WHERE chat_id = ? ORDER BY RANDOM() LIMIT 1`
	// err := d.db.QueryRowContext(d.cntxt, q).Scan(
	// 	&res.ChatID,
	// 	&res.Name,
	// 	&res.Description,
	// 	&res.Skills,
	// 	&res.YearsOfProgramming,
	// 	&res.LastCommand,
	// 	&res.IsImportant,
	// 	&res.Image)

	// if err == sql.ErrNoRows {
	// 	return LinusUser.User{}, storage.ErrNoSavedPages
	// }
	// if err != nil {
	// 	return LinusUser.User{}, e.Wrap("Can not get random user", err)
	// }
	return res, nil
}
