package sessions

import (
	"nubes/sum/db"
	"nubes/sum/utils"

	"gorm.io/gorm"
)

func Get(authenticationToken string) (*db.UserSession, error) {
	if token, err := utils.JwtVerify(authenticationToken); err == nil {
		return db.DB.UserSessions().FindById(token.JwtID())
	}

	// If the token fails JWT validation, return not found
	return nil, gorm.ErrRecordNotFound
}
