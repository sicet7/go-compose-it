package user

import (
	"errors"
	userType "github.com/sicet7/go-compose-it/src/database/entities/user"
	"gorm.io/gorm"
)

var (
	UsernameAlreadyInUse = errors.New("a user with the given username already exists")
)

type UserRepository struct {
	conn *gorm.DB
}

func NewRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		conn: db,
	}
}

func (r *UserRepository) FindByUsername(username string) (userType.User, error) {
	var user userType.User
	result := r.conn.First(&user, "username = ?", username)
	if result.Error != nil {
		return userType.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepository) Create(username string, password string) (userType.User, error) {
	var user userType.User
	_, err := r.FindByUsername(username)

	if err == nil {
		return userType.User{}, UsernameAlreadyInUse
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return userType.User{}, err
	}

	user, err = userType.New(username, password)
	if err != nil {
		return userType.User{}, err
	}

	result := r.conn.Create(&user)
	if result.Error != nil {
		return userType.User{}, result.Error
	}
	return user, nil
}
