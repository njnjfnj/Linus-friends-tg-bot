package storage

import (
	"LinusFriends/LinusUser"
	"errors"
)

var ErrNoFriends = errors.New("no friends :(")

const (
	SearchingByExperience = iota
	SearchingByRandom
	SearchingByLanguage
)

type Storage interface {
	AddNewUser(u LinusUser.User) error
	DeleteUser(chat_id int) error
	IsUserExists(chat_id int) (bool, error)
	GetUser(chat_id int) (LinusUser.User, error)
	UpdateUser(u LinusUser.User) error
	GetRandomUserForUser(chat_id int64, SearchingByWhat int, user LinusUser.User) (LinusUser.User, string, error)
}
