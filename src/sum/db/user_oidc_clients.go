package db

import (
	"errors"
	"nubes/sum/utils"

	"gorm.io/gorm"
)

type UserOidcClientRepository struct {
	handle *gorm.DB
}

func (db *Database) UserOidcClients() *UserOidcClientRepository {
	return &UserOidcClientRepository{db.handle}
}

func (r *UserOidcClientRepository) ListForUserID(userID string) []UserOidcClient {
	var oidcClients []UserOidcClient
	res := r.handle.Order("updated_at DESC").Find(&oidcClients)
	if res.Error != nil {
		panic(res.Error)
	}

	return oidcClients
}

func (r *UserOidcClientRepository) FindById(id string) (*UserOidcClient, error) {
	oidcClient := UserOidcClient{}

	res := r.handle.First(&oidcClient, "id = ?", id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &UserOidcClient{}, res.Error
	} else if res.Error != nil {
		panic(res.Error)
	}

	return &oidcClient, nil
}

func (r *UserOidcClientRepository) FindByUserAndClientID(userID, clientID string) (*UserOidcClient, error) {
	oidcClient := &UserOidcClient{}

	res := r.handle.First(oidcClient, "user_id = ? AND client_id = ?", userID, clientID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &UserOidcClient{}, res.Error
	} else if res.Error != nil {
		panic(res.Error)
	}

	return oidcClient, nil
}

func (r *UserOidcClientRepository) Create(UserOidcClient *UserOidcClient) {
	if err := r.handle.Create(&UserOidcClient).Error; err != nil {
		panic(err)
	}
}

func (r *UserOidcClientRepository) Delete(UserOidcClient *UserOidcClient) {
	if err := r.handle.Delete(&UserOidcClient).Error; err != nil {
		panic(err)
	}
}

func (r *UserOidcClientRepository) DeleteFor(ID string, userID string) error {
	result := r.handle.Delete(&UserOidcClient{}, "id = ? AND user_id = ?", ID, userID)
	if result.Error != nil {
		panic(result.Error)
	}

	if result.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

type UserOidcClient struct {
	Model
	UserID   string          `json:"-"`
	ClientID string          `json:"client_id"`
	Scopes   pipeStringArray `json:"scopes"`
}

func (u *UserOidcClient) AddScopes(scopes []string) {
	for _, v := range scopes {
		if !utils.Contains(u.Scopes, v) {
			u.Scopes = append(u.Scopes, v)
		}
	}
}

func (u *UserOidcClient) RevokeScope(scope string) {
	u.Scopes = utils.Select(u.Scopes, func(s string) bool {
		return s != scope
	})
}

func (u *UserOidcClient) HasScope(scope string) bool {
	return utils.Contains(u.Scopes, scope)
}

func (u *UserOidcClient) GetClaims() map[string]interface{} {
	result := make(map[string]interface{})
	user, err := DB.Users().FindById(u.UserID)
	if err != nil {
		panic(err)
	}

	for _, scope := range u.Scopes {
		switch scope {
		case "profile":
			result["preferred_username"] = user.Username
			result["updated_at"] = user.UpdatedAt.Unix()
			if user.Name != "" {
				result["name"] = user.Name
			}
			if user.Picture != "" {
				result["picture"] = user.Picture
			}
			if user.Birthdate != nil {
				result["birthdate"] = user.Birthdate.Format("2006-01-02")
			}
			if user.Zoneinfo != "" {
				result["zoneinfo"] = user.Zoneinfo
			}
			if user.Locale != "" {
				result["locale"] = user.Locale
			}
		case "email":
			if user.Email != "" {
				result["email"] = user.Email
				result["email_verified"] = user.EmailVerified
			}
		case "phone":
			if user.PhoneNumber != "" {
				result["phone_number"] = user.PhoneNumber
				result["phone_number_verified"] = user.PhoneNumberVerified
			}
		}
	}

	return result
}
