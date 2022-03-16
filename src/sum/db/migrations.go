package db

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&User{},

		&OidcClient{},
		&OidcClientRedirectURI{},
		&OidcClientResponseType{},
		&OidcClientGrantType{},
		&OidcClientContact{},
		&OidcClientLocalizedDetail{},
		&OidcClientDefaultACRValue{},
		&OidcClientRequestURI{},
	)
}
