package db

import (
	"errors"
	"log"
	"nubes/sum/utils"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	handle *gorm.DB
}

func (db *Database) Users() *UserRepository {
	return &UserRepository{db.handle}
}

type User struct {
	gorm.Model
	Username       string
	PasswordDigest string
	IsAdmin        bool
}

func (u *UserRepository) Count() int64 {
	var count int64
	u.handle.Model(&User{}).Count(&count)

	return count
}

func (u *UserRepository) FindById(id int) (*User, error) {
	user := User{}

	res := u.handle.First(&user, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &User{}, res.Error
	} else if res.Error != nil {
		log.Panicf("Could not load user: %v", res.Error)
	}

	return &user, nil
}

func (u UserRepository) FindByCredentials(identifier, password string) (*User, error) {
	user := User{}

	res := u.handle.First(&user, "username = ?", identifier)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &User{}, res.Error
	} else if res.Error != nil {
		log.Panicf("Could not load user: %v", res.Error)
	}

	if !user.VerifyPassword(password) {
		return &User{}, errors.New("Invalid username or password")
	}

	return &user, nil
}

func (u *UserRepository) Create(user *User) error {
	return u.handle.Create(&user).Error
}

func (u *UserRepository) Update(user *User) error {
	return u.handle.Save(&user).Error
}

func (u *UserRepository) Delete(user *User) error {
	return u.handle.Delete(&user).Error
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

func (u *User) SetPassword(password string) {
	u.PasswordDigest = utils.HashPassword(password)
}

func (u *User) VerifyPassword(password string) bool {
	return utils.VerifyPassword(password, u.PasswordDigest)
}
