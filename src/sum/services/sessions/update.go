package sessions

import (
	"errors"
	"log"
	"nubes/sum/db"
	"nubes/sum/utils"

	"github.com/lestrrat-go/jwx/jwt"
)

func Update(actor *db.User, session *db.UserSession, password string) (*db.UserSession, error) {
	if actor.ID != session.ID {
		return nil, errors.New("Cannot update someone else's session")
	}

	if !actor.VerifyPassword(password) {
		return nil, errors.New("Invalid username or password")
	}

	session, err := db.DB.UserSessions().Update(session)

	// Prepare an access token
	token, err := jwt.NewBuilder().
		JwtID(session.ID).
		Subject(actor.ID).
		Expiration(session.ExpiresAt).
		Audience([]string{"sessions"}).
		Build()
	if err != nil {
		log.Panicf("Could not build the JWT token: %v", err)
	}

	session.SignedToken = utils.JwtSign(token, "ES256")
	return session, nil
}
