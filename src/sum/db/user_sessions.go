package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserSessionRepository struct {
	handle *gorm.DB
}

func (db *Database) UserSessions() *UserSessionRepository {
	return &UserSessionRepository{db.handle}
}

func (r *UserSessionRepository) ListForUserID(userID, orderBy string) []*UserSession {
	if orderBy == "" {
		orderBy = "updated_at desc"
	}

	var sessions []*UserSession
	res := r.handle.Order(orderBy).Find(&sessions)
	if res.Error != nil {
		panic(res.Error)
	}

	return sessions
}

func (r *UserSessionRepository) FindById(id string) (*UserSession, error) {
	session := UserSession{}

	res := r.handle.First(&session, "id = ?", id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &UserSession{}, res.Error
	} else if res.Error != nil {
		panic(res.Error)
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

func (r *UserSessionRepository) DeleteFor(ID string, userID string) error {
	result := r.handle.Delete(&UserSession{}, "id = ? AND user_id = ?", ID, userID)
	if result.Error != nil {
		panic(result.Error)
	}

	if result.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserSessionRepository) CleanupExpired() {
	if err := r.handle.Where("expires_at < date('now')").Delete(&UserSession{}).Error; err != nil {
		panic(err)
	}
}

type UserSession struct {
	Model
	ExpiresAt time.Time `json:"-"`
	UserID    string    `json:"-"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`

	SignedToken string `gorm:"-" json:"-"`
}
