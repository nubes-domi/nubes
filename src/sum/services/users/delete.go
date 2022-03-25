package users

import (
	"nubes/sum/db"
	"nubes/sum/services"
)

func Delete(actor *db.User, userID string) error {
	if !actor.Admin {
		return &services.ForbiddenError{}
	}

	if actor.ID == userID {
		return &services.ConflictError{Detail: "cannot_self_delete"}
	}

	return db.DB.Users().Delete(userID)
}
