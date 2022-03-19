package db

import (
	"errors"
	"log"
	"nubes/sum/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (r *UserRepository) Count() int64 {
	var count int64
	r.handle.Model(&User{}).Count(&count)

	return count
}

func (r *UserRepository) FindById(id uint) (*User, error) {
	user := User{}

	res := r.handle.First(&user, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &User{}, res.Error
	} else if res.Error != nil {
		log.Panicf("Could not load user: %v", res.Error)
	}

	return &user, nil
}

func (r *UserRepository) FindByCredentials(identifier, password string) (*User, error) {
	user := User{}

	res := r.handle.First(&user, "username = ?", identifier)
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

func (r *UserRepository) Create(user *User) error {
	return r.handle.Create(&user).Error
}

func (r *UserRepository) Update(user *User) error {
	return r.handle.Save(&user).Error
}

func (r *UserRepository) Delete(user *User) error {
	return r.handle.Delete(&user).Error
}

func (u *User) SetPassword(password string) {
	u.PasswordDigest = utils.HashPassword(password)
}

func (u *User) VerifyPassword(password string) bool {
	return utils.VerifyPassword(password, u.PasswordDigest)
}

func (u *User) NewSession(c *gin.Context) UserSession {
	return UserSession{
		ID:        uuid.New().String(),
		UserID:    u.ID,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
		ExpiresAt: time.Now().Add(time.Duration(time.Now().Year()) * 10),
	}
}

type UserOidcScopes struct {
	gorm.Model
	UserID       uint
	OidcClientID string
	Scope        string
}

type UserOidcSession struct {
	gorm.Model
	UserID       uint
	OidcClientID string
	CodeDigest   string
}
