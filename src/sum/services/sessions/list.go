package sessions

import (
	"nubes/sum/db"
)

func List(actor *db.User, orderBy string) ([]*db.UserSession, error) {
	return db.DB.UserSessions().ListForUserID(actor.ID, orderBy), nil
}
