package users

import (
	"nubes/sum/db"
	"nubes/sum/services"
)

func Get(actor *db.User, userID string) (*db.User, error) {
	if actor.ID != userID && !actor.Admin {
		return nil, &services.ForbiddenError{}
	}

	return db.DB.Users().FindById(userID)
}
