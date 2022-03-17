package db

import (
	"bytes"
	"encoding/base64"
	"nubes/sum/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string
	PasswordDigest string
}

type UserSession struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UserID    int
	UserAgent string
	IPAddress string
}

type UserOidcScopes struct {
	gorm.Model
	UserID       int
	OidcClientID string
	Scope        string
}

type UserOidcSession struct {
	gorm.Model
	UserID       int
	OidcClientID string
	CodeDigest   string
}

func FindUserByUsername(username string) (User, bool) {
	user := User{}

	res := DB.First(&user, "username = ?", username)
	if res.Error != nil {
		return user, false
	}

	return user, true
}

func (u *User) SetPassword(password string) {
	u.PasswordDigest = utils.HashPassword(password)
}

func (u *User) VerifyPassword(password string) bool {
	pieces := strings.Split(u.PasswordDigest, "$")

	seed, _ := base64.RawStdEncoding.DecodeString(pieces[1])
	expected, _ := base64.RawStdEncoding.DecodeString(pieces[0])
	actual := utils.HashPasswordWithSeed(password, seed)

	return bytes.Equal(actual, expected)
}
