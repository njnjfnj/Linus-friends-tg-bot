package storage

import (
	"LinusFriends/LinusUser"
	"errors"
)

var ErrNoSavedPages = errors.New("no saved pages")

type Storage interface {
	AddNewUser(u LinusUser.User) error
	DeleteUser(chat_id int) error
	IsUserExists(chat_id int) (bool, error)
	GetUser(chat_id int) (LinusUser.User, error)
	UpdateUser(u LinusUser.User) error
	GetRandomUser() (LinusUser.User, error)
}
