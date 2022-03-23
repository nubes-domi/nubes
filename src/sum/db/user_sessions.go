package db

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserSessionRepository struct {
	handle *gorm.DB
}

func (db *Database) UserSessions() *UserSessionRepository {
	return &UserSessionRepository{db.handle}
}

func (r *UserSessionRepository) FindById(id string) (*UserSession, error) {
	session := UserSession{}

	res := r.handle.First(&session, "id = ?", id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &UserSession{}, res.Error
	} else if res.Error != nil {
		log.Panicf("Could not load user session: %v", res.Error)
	}

	return &session, nil
}

func (r *UserSessionRepository) Create(userSession *UserSession) {
	if err := r.handle.Create(&userSession).Error; err != nil {
		panic(err)
	}
}

func (r *UserSessionRepository) Delete(userSession *UserSession) {
	if err := r.handle.Delete(&userSession).Error; err != nil {
		panic(err)
	}
}

func (r *UserSessionRepository) CleanupExpired() {
	if err := r.handle.Where("expires_at < date('now')").Delete(&UserSession{}).Error; err != nil {
		panic(err)
	}
}

type UserSession struct {
	Model
	ExpiresAt time.Time
	UserID    string
	UserAgent string
	IPAddress string

	SignedToken string `gorm:"-"`
}
