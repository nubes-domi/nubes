package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string
	PasswordDigest string
}

type UserOidcScopes struct {
	gorm.Model
	UserID       string
	OidcClientID string
	Scope        string
}

type UserOidcSession struct {
	gorm.Model
	UserID       string
	OidcClientID string
	CodeDigest   string
}
