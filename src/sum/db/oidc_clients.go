package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"nubes/sum/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OidcClientRepository struct {
	handle *gorm.DB
}

func (db *Database) OidcClients() *OidcClientRepository {
	return &OidcClientRepository{db.handle}
}

func (r *OidcClientRepository) FindById(id string) (*OidcClient, error) {
	client := OidcClient{}

	res := r.handle.First(&client, "id = ?", id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &OidcClient{}, res.Error
	} else if res.Error != nil {
		log.Panicf("Could not load client: %v", res.Error)
	}

	client.ClientIDIssuedAt = client.CreatedAt.Unix()
	return &client, nil
}

func (r *OidcClientRepository) Create(client *OidcClient) error {
	err := r.handle.Create(&client).Error
	client.ClientIDIssuedAt = client.CreatedAt.Unix()
	return err
}

func (r *OidcClientRepository) Update(client *OidcClient) error {
	return r.handle.Save(&client).Error
}

func (r *OidcClientRepository) Delete(client *OidcClient) error {
	return r.handle.Delete(&client).Error
}

// Note, the ClientSecret is stored plaintext, not hashed.
// This is terrible, but the OpenID Connect Registration specification specifies
// that clients must be able to retrieve the ClientSecret through the registration_client_uri
type OidcClient struct {
	ID                                string          `json:"client_id" gorm:"primaryKey"`
	CreatedAt                         time.Time       `json:"-"`
	UpdatedAt                         time.Time       `json:"-"`
	ClientSecret                      string          `json:"client_secret"`
	ClientIDIssuedAt                  int64           `json:"client_id_issued_at" gorm:"-"`
	ClientSecretExpiresAt             int64           `json:"client_secret_expires_at" gorm:"-"`
	ApplicationType                   string          `json:"application_type"`
	JwksURI                           string          `json:"jwks_uri,omitempty"`
	Jwks                              jwkSet          `gorm:"type:text" json:"jwks,omitempty"`
	SectorIdentifierURI               string          `json:"sector_identifier_uri,omitempty"`
	SubjectType                       string          `json:"subject_type"`
	IDTokenSignedResponseAlg          string          `json:"id_token_signed_response_alg"`
	IDTokenEncryptedResponseAlg       string          `json:"id_token_encrypted_response_alg,omitempty"`
	IDTokenEncryptedResponseEnc       string          `json:"id_token_encrypted_response_enc,omitempty"`
	UserinfoSignedResponseAlg         string          `json:"userinfo_signed_response_alg,omitempty"`
	UserinfoEncryptedResponseAlg      string          `json:"userinfo_encrypted_response_alg,omitempty"`
	UserinfoEncryptedResponseEnc      string          `json:"userinfo_encrypted_response_enc,omitempty"`
	RequestObjectSignedResponseAlg    string          `json:"request_object_signed_response_alg,omitempty"`
	RequestObjectEncryptedResponseAlg string          `json:"request_object_encrypted_response_alg,omitempty"`
	RequestObjectEncryptedResponseEnc string          `json:"request_object_encrypted_response_enc,omitempty"`
	TokenEndpointAuthMethod           string          `json:"token_endpoint_auth_method,omitempty"`
	TokenEndpointAuthSigningAlg       string          `json:"token_endpoint_auth_signing_alg,omitempty"`
	DefaultMaxAge                     int             `json:"default_max_age,omitempty"`
	RequireAuthTime                   bool            `json:"require_auth_time,omitempty"`
	InitiateLoginURI                  string          `json:"initiate_login_uri,omitempty"`
	RedirectURIs                      pipeStringArray `gorm:"type:text" json:"redirect_uris"`
	ResponseTypes                     pipeStringArray `gorm:"type:text" json:"response_types,omitempty"`
	GrantTypes                        pipeStringArray `gorm:"type:text" json:"grant_types,omitempty"`
	Contacts                          pipeStringArray `gorm:"type:text" json:"contacts,omitempty"`
	DefaultACRValues                  pipeStringArray `gorm:"type:text" json:"default_acr_values,omitempty"`
	RequestURIs                       pipeStringArray `gorm:"type:text" json:"request_uris,omitempty"`

	// For things like client_uri#es
	LocalizedDetails []OidcClientLocalizedDetail `json:"-"`

	// For the client to make changes
	RegistrationAccessToken       string `json:"registration_access_token,omitempty" gorm:"-"`
	RegistrationAccessTokenDigest string `json:"-"`
	RegistrationClientURI         string `json:"registration_client_uri,omitempty"`
}

// Handles client_name, logo_uri, client_uri, policy_uri, tos_uri
type OidcClientLocalizedDetail struct {
	gorm.Model
	OidcClientID string
	Locale       string
	Field        string
	Value        string
}

func BuildOpenIDClient(c *gin.Context) OidcClient {
	client := OidcClient{
		ID:                      uuid.New().String(),
		ClientSecret:            utils.RandBase64(48),
		RegistrationAccessToken: utils.RandBase64(48),
	}
	client.RegistrationAccessTokenDigest = utils.Sha256String([]byte(client.RegistrationAccessToken))

	// To parse this request we are binding the body twice (to a well defined json structure and a map)
	// Gin however will prevent us from doing so.
	// So we will read the raw body bytes and parse them twice

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Panicf("Could not bind json: %v", err)
	}

	err = json.Unmarshal(reqBody, &client)
	if err != nil {
		log.Panicf("Could not bind json: %v", err)
	}

	request := make(map[string]interface{})
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		log.Panicf("Could not bind json: %v", err)
	}

	for k, v := range request {
		localised := strings.Split(k, "#")
		key := localised[0]

		switch key {
		case "client_name", "logo_uri", "client_uri", "policy_uri", "tos_uri":
			// these fields are potentially localised
			locale := ""
			if len(localised) > 1 {
				locale = localised[1]
			}

			client.LocalizedDetails = append(client.LocalizedDetails, OidcClientLocalizedDetail{
				Locale: locale,
				Field:  key,
				Value:  v.(string),
			})
		}
	}

	/*
	 * Fill in the blanks with defaults
	 */

	if len(client.ResponseTypes) == 0 {
		client.ResponseTypes = []string{"code"}
	}

	if len(client.GrantTypes) == 0 {
		client.GrantTypes = []string{"authorization_code"}
	}

	if client.ApplicationType == "" {
		client.ApplicationType = "web"
	}

	if client.IDTokenSignedResponseAlg == "" {
		client.IDTokenSignedResponseAlg = "RS256"
	}

	if client.TokenEndpointAuthMethod == "" {
		client.TokenEndpointAuthMethod = "client_secret_basic"
	}

	if client.SubjectType == "" {
		// We should probably prefer pairwise
		client.SubjectType = "public"
	}

	return client
}

func (client *OidcClient) getLocalizedDetail(field, locale string) string {
	row := OidcClientLocalizedDetail{}

	res := DB.handle.First(&row, "oidc_client_id = ? AND field = ? AND locale = ?", client.ID, field, locale)
	if res.Error == nil {
		return row.Value
	} else if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Panicf("Could not retrieve OIDC Client localized detail: %v", res.Error)
	}

	res = DB.handle.First(&row, "oidc_client_id = ? AND field = ? AND locale = ?", client.ID, field, "")
	if res.Error == nil {
		return row.Value
	} else if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Panicf("Could not retrieve OIDC Client localized detail: %v", res.Error)
	}

	return ""
}

func (client *OidcClient) GetClientName(locale string) string {
	return client.getLocalizedDetail("client_name", locale)
}

func (client *OidcClient) GetClientURI(locale string) string {
	return client.getLocalizedDetail("client_uri", locale)
}

func (client *OidcClient) GetLogoURI(locale string) string {
	return client.getLocalizedDetail("logo_uri", locale)
}

func (client *OidcClient) GetPolicyURI(locale string) string {
	return client.getLocalizedDetail("policy_uri", locale)
}

func (client *OidcClient) GetTosURI(locale string) string {
	return client.getLocalizedDetail("tos_uri", locale)
}

func (client *OidcClient) VerifyRegistrationToken(token string) bool {
	return utils.Sha256String([]byte(token)) == client.RegistrationAccessTokenDigest
}

func validateURI(uri string, allowFragments bool) error {
	parsed, err := url.Parse(uri)
	if err != nil {
		return err
	} else if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("URI uses invalid scheme")
	}

	if !allowFragments && parsed.Fragment != "" {
		return errors.New("URI has fragment")
	}

	return err
}

func validateURIArray(uris []string, allowFragments bool) error {
	for _, uri := range uris {
		if err := validateURI(uri, allowFragments); err != nil {
			return err
		}
	}

	return nil
}

func LoadClient(clientId string) OidcClient {
	client := OidcClient{}

	res := DB.handle.Preload(clause.Associations).First(&client, "id = ?", clientId)
	if res.Error != nil {
		log.Panicf("Error while Loading client: %v", res.Error)
	}

	return client
}

func (client *OidcClient) ExplodedResponseTypes() []string {
	result := []string{}
	for _, v := range client.ResponseTypes {
		types := strings.Split(v, " ")
		for _, t := range types {
			if !utils.Contains(result, t) {
				result = append(result, t)
			}
		}
	}

	return result
}

func (client *OidcClient) Validate() error {
	if err := validateURIArray(client.RedirectURIs, false); err != nil {
		return errors.New("invalid_redirect_uri")
	}

	grantTypes := client.GrantTypes

	responseGrantsDependencies := map[string][]string{
		"code":     {"authorization_code"},
		"id_token": {"implicit"},
		"token":    {"implicit"},
	}

	for _, responseType := range client.ExplodedResponseTypes() {
		switch responseType {
		case "code", "id_token", "token":
		default:
			return fmt.Errorf("Invalid response type %s", responseType)
		}

		for _, grant := range responseGrantsDependencies[responseType] {
			if !utils.Contains(grantTypes, grant) {
				return fmt.Errorf("\"%s\" response type requires \"%s\" grant which is not being registered", responseType, grant)
			}
		}
	}

	for _, grantType := range grantTypes {
		switch grantType {
		case "authorization_code", "implicit", "refresh_token":
		default:
			return fmt.Errorf("Invalid grant type %s", grantType)
		}
	}

	switch client.ApplicationType {
	case "web":
		for _, uri := range client.RedirectURIs {
			parsed, _ := url.Parse(uri)
			if parsed.Scheme != "https" || parsed.Host == "localhost" {
				return fmt.Errorf("Cannot use %s as Redirect URI for a web client (must be https://<notlocalhost>/<path>)", uri)
			}
		}
	case "native":
		for _, uri := range client.RedirectURIs {
			parsed, _ := url.Parse(uri)
			if parsed.Scheme != "http" || parsed.Host != "localhost" {
				return fmt.Errorf("Cannot use %s as Redirect URI for a native client (must be http://localhost/<path>)", uri)
			}
		}
	default:
		return fmt.Errorf("Invalid application type: %s. Permitted values are \"web\" and \"native\"", client.ApplicationType)
	}

	switch client.SubjectType {
	case "pairwise", "public":
	default:
		return fmt.Errorf("Invalid subject type: %s. Permitted values are \"public\" and \"pairwise\"", client.SubjectType)
	}

	if client.JwksURI != "" {
		if err := validateURI(client.JwksURI, true); err != nil {
			return err
		}

		if client.Jwks.Set != nil {
			return fmt.Errorf("RP cannot provide both jwks_uri and jwks")
		}
	}

	if client.InitiateLoginURI != "" {
		if err := validateURI(client.InitiateLoginURI, true); err != nil {
			return err
		}
	}

	if err := validateURIArray(client.RequestURIs, true); err != nil {
		return err
	}

	return nil
}
