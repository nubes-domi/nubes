package db

import (
	"errors"
	"nubes/sum/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
	handle *gorm.DB
}

func (db *Database) Users() *UserRepository {
	return &UserRepository{db.handle}
}

type User struct {
	Model
	Username       string `json:"username"`
	PasswordDigest string `json:"-"`
	IsAdmin        bool   `json:"is_admin" binding:"-"`

	Name                string   `json:"name,omitempty"`
	Picture             string   `json:"picture,omitempty"`
	Email               string   `json:"email,omitempty"`
	EmailVerified       bool     `json:"email_verified" binding:"-"`
	Birthdate           JSONDate `json:"birthdate,omitempty"`
	Zoneinfo            string   `json:"zoneinfo,omitempty"`
	PhoneNumber         string   `json:"phone_number,omitempty"`
	PhoneNumberVerified bool     `json:"phone_number_verified" binding:"-"`
}

func (r *UserRepository) New() *User {
	return &User{
		Model: Model{ID: GenID("usr")},
	}
}

func (r *UserRepository) Count() int64 {
	var count int64
	r.handle.Model(&User{}).Count(&count)

	return count
}

func (r *UserRepository) FindById(id string) (*User, error) {
	user := User{}

	res := r.handle.First(&user, "id = ?", id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &User{}, res.Error
	} else if res.Error != nil {
		panic(res.Error)
	}

	return &user, nil
}

func (r *UserRepository) FindByCredentials(identifier, password string) (*User, error) {
	user := User{}

	res := r.handle.First(&user, "username = ?", identifier)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &User{}, res.Error
	} else if res.Error != nil {
		panic(res.Error)
	}

	if !user.VerifyPassword(password) {
		return &User{}, errors.New("Invalid username or password")
	}

	return &user, nil
}

func (r *UserRepository) List(order string) []User {
	var users []User
	res := r.handle.Order(order).Find(&users)
	if res.Error != nil {
		panic(res.Error)
	}

	return users
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

func (u *User) NewSession(c *gin.Context) *UserSession {
	return &UserSession{
		Model:     Model{ID: GenID("usr_sess")},
		UserID:    u.ID,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
		ExpiresAt: time.Now().Add(time.Duration(time.Now().Year()) * 10),
	}
}

type UserOidcScopes struct {
	Model
	UserID       uint
	OidcClientID string
	Scope        string
}

type UserOidcSession struct {
	Model
	UserID       uint
	OidcClientID string
	CodeDigest   string
}
