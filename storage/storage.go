package storage

import (
	"LinusFriends/LinusUser"
	"LinusFriends/advertisement"
	"errors"
)

var ErrNoFriends = errors.New("no friends :(")
var ErrNoAds = errors.New("no adverts")

const (
	SearchingByExperience = iota
	SearchingByRandom
	SearchingByLanguage
)

type Storage interface {
	UserCount() (int, error)
	AddNewUser(u LinusUser.User) error
	DeleteUser(chat_id int) error
	IsUserExists(chat_id int) (bool, error)
	GetUser(chat_id int) (LinusUser.User, error)
	UpdateUser(u LinusUser.User) error
	GetRandomUserForUser(chat_id int64, SearchingByWhat int, user LinusUser.User) (LinusUser.User, string, error)
	DeleteSkill(chat_id int, skill string) error

	AddMatch(chat_id int64, user LinusUser.User) error
	GetMatches(chat_id int64) (string, error)
	SetMatches(chat_id int64, matchesLeft string) error

	AddNewAd(ad advertisement.Ad) error
	DeleteAd(advert_id int) error
	UpdateAdContent(content []byte, description string, advert_id int) error
	UpdateAdRatingAndViews(rating float32, seen int, rated int, advert_id int) error
	GetAd(advert_id int) (advertisement.Ad, error)
	GetAds() ([]advertisement.Ad, error)
	GetAdsIds() ([]int, error)
	IsAdExists(advert_id int) (bool, error)
}
