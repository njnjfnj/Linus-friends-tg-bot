package sqlite

import (
	LinusUser "LinusFriends/User"
	"LinusFriends/libs/e"
	"context"
	"database/sql"
	"log"

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
	q := `CREATE TABLE IF NOT EXISTS pages (chat_id INT, name TEXT, description TEXT, last_CMD TEXT, last_CMD_Pos INT, photo BLOB NOT NULL)`
	if _, err := d.db.ExecContext(context, q); err != nil {
		return e.Wrap("Can not create DB", err)
	}
	d.cntxt = context

	return nil
}

func (d *Database) AddNewUser(u LinusUser.User) {
	q := `INSERT INTO pages (chat_id, name, description, photo) VALUES(?, ?, ?, ?, ?, ?)`
	if _, err := d.db.ExecContext(d.cntxt, q, u.ChatID, u.Name, u.Description, u.LastCommand, u.Image); err != nil {
		log.Print("Can not save page", err)
	}
}

//func (d *Database)
