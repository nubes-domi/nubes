package sessions

import (
	"nubes/sum/db"
)

func Delete(actor *db.User, sessionID string) {
	db.DB.UserSessions().DeleteFor(sessionID, actor.ID)
}
