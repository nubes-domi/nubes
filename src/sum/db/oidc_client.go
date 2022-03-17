package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils"
)

type SimpleStringArray []string

func (n *SimpleStringArray) Scan(value interface{}) error {
	*n = strings.Split(string(value.(string)), "|")
	return nil
}

func (n *SimpleStringArray) Value() (driver.Value, error) {
	return driver.Value(strings.Join(*n, "|")), nil
}

type OidcClient struct {
	ID                                string                      `gorm:"primaryKey"`
	CreatedAt                         time.Time                   `json:"-"`
	UpdatedAt                         time.Time                   `json:"-"`
	ClientSecretDigest                string                      `json:"-"`
	ClientSecret                      string                      `json:"client_secret" gorm:"-"`
	ApplicationType                   string                      `json:"application_type"`
	JwksURI                           string                      `json:"jwks_uri"`
	Jwks                              string                      `json:"-"`
	SectorIdentifierURI               string                      `json:"sector_identifier_uri"`
	SubjectType                       string                      `json:"subject_type"`
	IDTokenSignedResponseAlg          string                      `json:"id_token_signed_response_alg"`
	IDTokenEncryptedResponseAlg       string                      `json:"id_token_encrypted_response_alg"`
	IDTokenEncryptedResponseEnc       string                      `json:"id_token_encrypted_response_enc"`
	UserinfoSignedResponseAlg         string                      `json:"userinfo_signed_response_alg"`
	UserinfoEncryptedResponseAlg      string                      `json:"userinfo_encrypted_response_alg"`
	UserinfoEncryptedResponseEnc      string                      `json:"userinfo_encrypted_response_enc"`
	RequestObjectSignedResponseAlg    string                      `json:"request_object_signed_response_alg"`
	RequestObjectEncryptedResponseAlg string                      `json:"request_object_encrypted_response_alg"`
	RequestObjectEncryptedResponseEnc string                      `json:"request_object_encrypted_response_enc"`
	TokenEndpointAuthMethod           string                      `json:"token_endpoint_auth_method"`
	TokenEndpointAuthSigningAlg       string                      `json:"token_endpoint_auth_signing_alg"`
	DefaultMaxAge                     int                         `json:"default_max_age"`
	RequireAuthTime                   bool                        `json:"require_auth_time"`
	InitiateLoginURI                  string                      `json:"initiate_login_uri"`
	RedirectURIs                      SimpleStringArray           `gorm:"type:text" json:"redirect_uris"`
	ResponseTypes                     SimpleStringArray           `gorm:"type:text" json:"response_types"`
	GrantTypes                        SimpleStringArray           `gorm:"type:text" json:"grant_types"`
	Contacts                          SimpleStringArray           `gorm:"type:text" json:"contacts"`
	DefaultACRValues                  SimpleStringArray           `gorm:"type:text" json:"default_acr_values"`
	RequestURIs                       SimpleStringArray           `gorm:"type:text" json:"request_uris"`
	LocalizedDetails                  []OidcClientLocalizedDetail `json:"-"`
}

// Handles client_name, logo_uri, client_uri, policy_uri, tos_uri
type OidcClientLocalizedDetail struct {
	gorm.Model
	OidcClientID string
	Locale       string
	Field        string
	Value        string
}

func toStringArray(field interface{}) []string {
	asArray := field.([]interface{})
	asStrings := []string{}

	for _, s := range asArray {
		asStrings = append(asStrings, s.(string))
	}

	return asStrings
}

func BuildOpenIDClient(c *gin.Context) OidcClient {
	client := OidcClient{}
	client.ID = uuid.New().String()

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
		case "jwks":
			// value is a nested object, make sure it is a valid JWK and store it as a JSON string
			plain, _ := json.Marshal(v)
			_, err = jwk.Parse(plain)
			if err != nil {
				log.Panicf("Invalid JWK object")
			}

			client.Jwks = string(plain)

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

	tx := DB.First(&row, "oidc_client_id = ? AND field = ? AND locale = ?", client.ID, field, locale)
	if tx.Error == nil {
		return row.Value
	}

	tx = DB.First(&row, "oidc_client_id = ? AND field = ? AND locale = ?", client.ID, field, "")
	if tx.Error == nil {
		return row.Value
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

func validateURI(uri string) error {
	_, err := url.ParseRequestURI(uri)
	return err
}

func validateURIArray(uris []string) error {
	for _, uri := range uris {
		if err := validateURI(uri); err != nil {
			return err
		}
	}

	return nil
}

func LoadClient(clientId string) OidcClient {
	client := OidcClient{}

	res := DB.Preload(clause.Associations).First(&client, "id = ?", clientId)
	if res.Error != nil {
		log.Panicf("Error while Loading client: %v", res.Error)
	}

	return client
}

func (client *OidcClient) Validate() error {
	redirectURIs := client.RedirectURIs
	if err := validateURIArray(redirectURIs); err != nil {
		return err
	}

	grantTypes := client.GrantTypes

	responseGrantsDependencies := map[string][]string{
		"code":                {"authorization_code"},
		"id_token":            {"implicit"},
		"token id_token":      {"implicit"},
		"code id_token":       {"authorization_code", "implicit"},
		"code token":          {"authorization_code", "implicit"},
		"code token id_token": {"authorization_code", "implicit"},
	}

	for _, responseType := range client.ResponseTypes {
		switch responseType {
		case "code", "id_token", "token id_token", "code id_token", "code token", "code token id_token":
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
		for _, uri := range redirectURIs {
			parsed, _ := url.Parse(uri)
			if parsed.Scheme != "https" || parsed.Host == "localhost" {
				return fmt.Errorf("Cannot use %s as Redirect URI for a web client (must be https://<notlocalhost>/<path>)", uri)
			}
		}
	case "native":
		for _, uri := range redirectURIs {
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
		if err := validateURI(client.JwksURI); err != nil {
			return err
		}

		if len(client.Jwks) != 0 {
			return fmt.Errorf("RP cannot provide both jwks_uri and jwks")
		}
	}

	if client.InitiateLoginURI != "" {
		if err := validateURI(client.InitiateLoginURI); err != nil {
			return err
		}
	}

	if err := validateURIArray(client.RequestURIs); err != nil {
		return err
	}

	return nil
}
