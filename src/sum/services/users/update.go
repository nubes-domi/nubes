package users

import (
	"nubes/sum/db"
	"nubes/sum/services"
	"nubes/sum/utils"
)

func Update(actor *db.User, updated *db.User) (*db.User, error) {
	if actor.ID != updated.ID && actor.Admin {
		return nil, &services.ForbiddenError{}
	}

	user, err := Get(actor, updated.ID)
	if err != nil {
		return nil, err
	}

	// Admins can promote and demote other users, but not themselves
	if actor.Admin && user.ID != actor.ID {
		user.Admin = updated.Admin
	}

	if user.Password != "" {
		user.PasswordDigest = utils.HashPassword(user.Password)
	}

	if updated.Email != "" && user.Email != updated.Email {
		user.Email = updated.Email
		user.EmailVerified = false
	}

	if updated.PhoneNumber != "" && user.PhoneNumber != updated.PhoneNumber {
		user.PhoneNumber = updated.PhoneNumber
		user.PhoneNumberVerified = false
	}

	if updated.Username != "" {
		user.Username = updated.Username
	}

	if updated.Name != "" {
		user.Name = updated.Name
	}

	if updated.Picture != "" {
		user.Picture = updated.Picture
	}

	if updated.Birthdate != nil {
		user.Birthdate = updated.Birthdate
	}

	if updated.Locale != "" {
		user.Locale = updated.Locale
	}

	if updated.Zoneinfo != "" {
		user.Zoneinfo = updated.Zoneinfo
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	return user, db.DB.Users().Update(user)
}
