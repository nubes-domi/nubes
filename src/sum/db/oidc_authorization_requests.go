package db

import (
	"errors"
	"log"
	"nubes/sum/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OidcAuthorizationRequestRepository struct {
	handle *gorm.DB
}

func (db *Database) OidcAuthorizationRequests() *OidcAuthorizationRequestRepository {
	return &OidcAuthorizationRequestRepository{db.handle}
}

func (r *OidcAuthorizationRequestRepository) FindByIdAndStage(id string, stage string) (*OidcAuthorizationRequest, error) {
	client := OidcAuthorizationRequest{}

	res := r.handle.First(&client, "id = ? AND stage = ?", id, stage)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &OidcAuthorizationRequest{}, res.Error
	} else if res.Error != nil {
		log.Panicf("Could not load client: %v", res.Error)
	}

	return &client, nil
}

func (r *OidcAuthorizationRequestRepository) Create(client *OidcAuthorizationRequest) error {
	return r.handle.Create(&client).Error
}

func (r *OidcAuthorizationRequestRepository) Update(client *OidcAuthorizationRequest) error {
	return r.handle.Save(&client).Error
}

func (r *OidcAuthorizationRequestRepository) Delete(client *OidcAuthorizationRequest) error {
	return r.handle.Delete(&client).Error
}

type OidcAuthorizationRequest struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// authorization, code
	Stage string
	// Stored after consent is granted
	UserID uint

	ClientID     string `form:"client_id"`
	Scope        string `form:"scope"`
	ResponseType string `form:"response_type"`
	RedirectURI  string `form:"redirect_uri"`
	State        string `form:"state"`
	ResponseMode string `form:"response_mode"`
	Nonce        string `form:"nonce"`
	Display      string `form:"display"`
	Prompt       string `form:"prompt"`
	MaxAge       uint   `form:"max_age"`
	UILocales    string `form:"ui_locales"`
	IDTokenHint  string `form:"id_token_hint"`
	LoginHint    string `form:"login_hint"`
	ACRValues    string `form:"acr_values"`

	// Form binding only
	RequestURI string `form:"request_uri" gorm:"-"`
	Request    string `form:"request" gorm:"-"`
}

func (r *OidcAuthorizationRequest) HasResponseType(responseType string) bool {
	types := strings.Split(r.ResponseType, " ")
	return utils.Contains(types, responseType)
}

func (r *OidcAuthorizationRequest) GetResponseMode() string {
	if r.ResponseMode != "" {
		return r.ResponseMode
	}

	if r.ResponseType == "code" {
		return "query"
	}

	return "fragment"
}
