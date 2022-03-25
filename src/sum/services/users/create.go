package users

import (
	"nubes/sum/db"
	"nubes/sum/services"
	"nubes/sum/utils"
)

func Create(actor *db.User, user *db.User) (*db.User, error) {
	if !actor.Admin {
		return nil, &services.ForbiddenError{}
	}

	user.ID = db.GenID("usr")
	user.EmailVerified = false
	user.PhoneNumberVerified = false

	if user.Password != "" {
		user.PasswordDigest = utils.HashPassword(user.Password)
	}

	err := user.Validate()
	if err != nil {
		return nil, err
	}

	err = db.DB.Users().Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
