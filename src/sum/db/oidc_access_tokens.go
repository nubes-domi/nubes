package db

import (
	"errors"
	"log"
	"nubes/sum/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OidcAccessTokensRepository struct {
	handle *gorm.DB
}

func (db *Database) OidcAccessTokens() *OidcAccessTokensRepository {
	return &OidcAccessTokensRepository{db.handle}
}

func (r *OidcAccessTokensRepository) Find(accessToken string) (*OidcAccessToken, error) {
	obj := OidcAccessToken{}

	parts := strings.Split(accessToken, ":")
	if len(parts) != 2 {
		return &OidcAccessToken{}, gorm.ErrRecordNotFound
	}

	res := r.handle.First(&obj, "id = ?", parts[0])
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &OidcAccessToken{}, res.Error
	} else if res.Error != nil {
		log.Panicf("Could not load client: %v", res.Error)
	}

	if obj.SecretDigest != utils.Sha256String([]byte(parts[1])) {
		return &OidcAccessToken{}, gorm.ErrRecordNotFound
	}

	return &obj, nil
}

func (r *OidcAccessTokensRepository) Create(client *OidcAccessToken) error {
	return r.handle.Create(&client).Error
}

func (r *OidcAccessTokensRepository) Update(client *OidcAccessToken) error {
	return r.handle.Save(&client).Error
}

func (r *OidcAccessTokensRepository) Delete(client *OidcAccessToken) error {
	return r.handle.Delete(&client).Error
}

type OidcAccessToken struct {
	ID           string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ExpiresAt    time.Time
	SecretDigest string
	ClientID     string
	UserID       uint
	Scope        pipeStringArray
}
