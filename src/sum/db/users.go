package db

import (
	"errors"
	"nubes/sum/utils"
	"regexp"
	"strings"
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
	Password       string `json:"password,omitempty" gorm:"-"`
	PasswordDigest string `json:"-"`
	Admin          bool   `json:"admin" binding:"-"`

	Name                string    `json:"name,omitempty"`
	Picture             string    `json:"picture,omitempty"`
	Email               string    `json:"email,omitempty"`
	EmailVerified       bool      `json:"email_verified" binding:"-"`
	Birthdate           *JSONDate `json:"birthdate,omitempty"`
	Locale              string    `json:"locale,omitempty"`
	Zoneinfo            string    `json:"zoneinfo,omitempty"`
	PhoneNumber         string    `json:"phone_number,omitempty"`
	PhoneNumberVerified bool      `json:"phone_number_verified" binding:"-"`
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

func (r *UserRepository) FindByIdentifier(identifier string) (*User, error) {
	if identifier == "" {
		return nil, gorm.ErrRecordNotFound
	}

	user := User{}

	identifier = strings.ToLower(identifier)
	res := r.handle.First(&user, "lower(username) = ? OR lower(email) = ?", identifier, identifier)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &User{}, res.Error
	} else if res.Error != nil {
		panic(res.Error)
	}

	return &user, nil
}

func (r *UserRepository) FindByCredentials(identifier, password string) (*User, error) {
	user, err := r.FindByIdentifier(identifier)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, err
	} else if err != nil {
		panic(err)
	}

	if !user.VerifyPassword(password) {
		return &User{}, errors.New("Invalid username or password")
	}

	return user, nil
}

func (r *UserRepository) List(orderBy string) []*User {
	var users []*User
	res := r.handle.Order(orderBy).Find(&users)
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

func (r *UserRepository) Delete(userID string) error {
	return r.handle.Delete("id = ?", userID).Error
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

func (u *User) GrantedScopesForClient(clientID string) []string {
	userOidcClient, err := DB.UserOidcClients().FindByUserAndClientID(u.ID, clientID)
	if err != nil {
		return []string{}
	}

	return userOidcClient.Scopes
}

func (u *User) Validate() error {
	if u.Password != "" {
		if len(u.Password) < 8 {
			return &ValidationError{Field: "password", Detail: "too_short"}
		}
		u.PasswordDigest = utils.HashPassword(u.Password)
	}

	if u.Email != "" {
		if m, _ := regexp.MatchString(`\A.+@.+\z`, u.Email); !m {
			return &ValidationError{Field: "email"}
		}
	}

	if u.PhoneNumber != "" {
		if m, _ := regexp.MatchString(`\A+?[\d ]{5,20}\z`, u.PhoneNumber); !m {
			return &ValidationError{Field: "phone_number"}
		}
	}

	if u.Birthdate != nil && u.Birthdate.After(time.Now()) {
		return &ValidationError{Field: "birthdate"}
	}

	return nil
}
