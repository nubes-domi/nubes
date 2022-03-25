package users

import (
	"nubes/sum/db"
	"nubes/sum/services"
)

func List(actor *db.User, orderBy string) ([]*db.User, error) {
	if !actor.Admin {
		return nil, &services.ForbiddenError{}
	}

	return db.DB.Users().List(orderBy), nil
}
