package user

import (
	"github.com/google/uuid"
	passwordHash "github.com/sicet7/go-compose-it/src/password"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID         `gorm:"primaryKey;type:string;size:36;<-:create" json:"id"`
	Username  string            `gorm:"unique" json:"username"`
	Password  passwordHash.Hash `gorm:"type:string;size:200;serializer:password" json:"-"`
	CreatedAt time.Time         `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *User) RequestTag() string {
	return u.ID.String()
}

func (u *User) VerifyPassword(password string) bool {
	return u.Password.VerifyPassword(password)
}

func New(username string, password string) (User, error) {
	hash, err := passwordHash.Argon2idFromPassword(password, passwordHash.DefaultArgon2Params())
	if err != nil {
		return User{}, err
	}
	return User{
		ID:       uuid.New(),
		Username: username,
		Password: &hash,
	}, nil
}
