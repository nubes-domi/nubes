package db

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&User{},
		&UserSession{},

		&OidcClient{},
		&OidcClientLocalizedDetail{},
		&OidcAuthorizationRequest{},
	)
}
