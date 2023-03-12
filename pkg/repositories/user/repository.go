package user

import (
	"errors"
	"gorm.io/gorm"
)

var (
	UsernameAlreadyInUse = errors.New("a user with the given username already exists")
)

type UserRepository struct {
	conn *gorm.DB
}

func (r *UserRepository) FindUserByUsername(username string) (User, error) {
	var user User
	result := r.conn.First(&user, "username = ?", username)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (r *UserRepository) CreateUser(username string, password string) (User, error) {
	var user User
	_, err := r.FindUserByUsername(username)

	if err == nil {
		return User{}, UsernameAlreadyInUse
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, err
	}

	user, err = newUser(username, password)
	if err != nil {
		return User{}, err
	}

	result := r.conn.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}
