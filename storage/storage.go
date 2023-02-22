package storage

import (
	usr "LinusFriends/User"
)

type Storage interface {
	AddNewUser(u usr.User) error
	DeleteUser(chat_id int) error
	IsUserExists(u usr.User) (bool, error)
	GetUser(chat_id int) (usr.User, error)
	UpdateUser(u usr.User) error
	GetRandomUser() error
}
