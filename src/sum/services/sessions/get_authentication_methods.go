package sessions

import (
	"nubes/sum/db"
)

func GetAuthenticationMethods(identifier string) (*db.User, []string, error) {
	user, err := db.DB.Users().FindByIdentifier(identifier)
	if err != nil {
		return nil, []string{}, err
	}

	// This is a stub, later on it support Webauthn, Email, etc
	return user, []string{"password"}, nil
}
