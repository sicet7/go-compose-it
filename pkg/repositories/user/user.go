package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:string;size:36;<-:create" json:"id"`
	Username  string         `gorm:"unique" json:"username"`
	Password  PasswordHash   `gorm:"type:string;size:200" json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *User) RequestTag() string {
	return u.ID.String()
}

func (u *User) VerifyPassword(password string) bool {
	result, err := u.Password.ComparePlainText(password)
	return err == nil && result
}

func newUser(username string, password string) (User, error) {
	passwordHash, err := newPasswordHash(password)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:       uuid.New(),
		Username: username,
		Password: passwordHash,
	}, nil
}
