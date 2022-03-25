package sessions

import (
	"errors"
	"log"
	"nubes/sum/db"
	"nubes/sum/utils"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
)

func Create(username, password, userAgent, ipAddress string) (*db.UserSession, error) {
	user, err := db.DB.Users().FindByCredentials(username, password)
	if err != nil {
		return nil, errors.New("invalid_credentials")
	} else {
		// Genearate and save the new session
		session := &db.UserSession{
			Model:     db.Model{ID: db.GenID("usr_sess")},
			UserID:    user.ID,
			UserAgent: userAgent,
			IPAddress: ipAddress,
			ExpiresAt: time.Now().Add(time.Duration(time.Now().Year()) * 10),
		}
		db.DB.UserSessions().Create(session)

		// Prepare an access token
		token, err := jwt.NewBuilder().
			JwtID(session.ID).
			Subject(user.ID).
			Expiration(session.ExpiresAt).
			Audience([]string{"sessions"}).
			Build()
		if err != nil {
			log.Panicf("Could not build the JWT token: %v", err)
		}

		session.SignedToken = utils.JwtSign(token, "RS256")
		return session, nil
	}
}
