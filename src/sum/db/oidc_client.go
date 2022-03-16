package db

import (
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
	"gorm.io/gorm/utils"
)

type OidcClient struct {
	ID                                string `gorm:"primaryKey"`
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
	ClientSecretDigest                string
	ApplicationType                   string `json:"application_type"`
	JWKsURI                           string `json:"jwks_uri"`
	JWKs                              string `json:"-"`
	SectorIdentifierURI               string `json:"sector_identifier_uri"`
	SubjectType                       string `json:"subject_type"`
	IDTokenSignedResponseAlg          string `json:"id_token_signed_response_alg"`
	IDTokenEncryptedResponseAlg       string `json:"id_token_encrypted_response_alg"`
	IDTokenEncryptedResponseEnc       string `json:"id_token_encrypted_response_enc"`
	UserinfoSignedResponseAlg         string `json:"userinfo_signed_response_alg"`
	UserinfoEncryptedResponseAlg      string `json:"userinfo_encrypted_response_alg"`
	UserinfoEncryptedResponseEnc      string `json:"userinfo_encrypted_response_enc"`
	RequestObjectSignedResponseAlg    string `json:"request_object_signed_response_alg"`
	RequestObjectEncryptedResponseAlg string `json:"request_object_encrypted_response_alg"`
	RequestObjectEncryptedResponseEnc string `json:"request_object_encrypted_response_enc"`
	TokenEndpointAuthMethod           string `json:"token_endpoint_auth_method"`
	TokenEndpointAuthSigningAlg       string `json:"token_endpoint_auth_signing_alg"`
	DefaultMaxAge                     int    `json:"default_max_age"`
	RequireAuthTime                   bool   `json:"require_auth_time"`
	InitiateLoginURI                  string `json:"initiate_login_uri"`

	RedirectURIs     []OidcClientRedirectURI
	ResponseTypes    []OidcClientResponseType
	GrantTypes       []OidcClientGrantType
	Contacts         []OidcClientContact `json:"-"`
	LocalizedDetails []OidcClientLocalizedDetail
	DefaultACRValues []OidcClientDefaultACRValue
	RequestURIs      []OidcClientRequestURI
}

type OidcClientRedirectURI struct {
	gorm.Model
	OidcClientID string
	RedirectURI  string
}

type OidcClientResponseType struct {
	gorm.Model
	OidcClientID string
	ResponseType string
}

type OidcClientGrantType struct {
	gorm.Model
	OidcClientID string
	GrantType    string
}

type OidcClientContact struct {
	gorm.Model
	OidcClientID string
	Contact      string
}

// Handles client_name, logo_uri, client_uri, policy_uri, tos_uri
type OidcClientLocalizedDetail struct {
	gorm.Model
	OidcClientID string
	Locale       string
	Field        string
	Value        string
}

type OidcClientDefaultACRValue struct {
	gorm.Model
	OidcClientID string
	ACRValue     string
}

type OidcClientRequestURI struct {
	gorm.Model
	OidcClientID string
	RequestURI   string
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

			client.JWKs = string(plain)

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
		case "redirect_uris":
			uris := toStringArray(v)
			for _, uri := range uris {
				client.RedirectURIs = append(client.RedirectURIs, OidcClientRedirectURI{
					RedirectURI: uri,
				})
			}
		case "response_types":
			response_types := toStringArray(v)
			for _, response_type := range response_types {
				client.ResponseTypes = append(client.ResponseTypes, OidcClientResponseType{
					ResponseType: response_type,
				})
			}
		case "grant_types":
			grant_types := toStringArray(v)
			for _, grant_type := range grant_types {
				client.GrantTypes = append(client.GrantTypes, OidcClientGrantType{
					GrantType: grant_type,
				})
			}
		case "contacts":
			contacts := toStringArray(v)
			for _, contact := range contacts {
				client.Contacts = append(client.Contacts, OidcClientContact{
					Contact: contact,
				})
			}
		case "default_acr_values":
			acr_values := toStringArray(v)
			for _, acr_value := range acr_values {
				client.DefaultACRValues = append(client.DefaultACRValues, OidcClientDefaultACRValue{
					ACRValue: acr_value,
				})
			}
		case "request_uris":
			uris := toStringArray(v)
			for _, uri := range uris {
				client.RequestURIs = append(client.RequestURIs, OidcClientRequestURI{
					RequestURI: uri,
				})
			}
		}
	}

	/*
	 * Fill in the blanks with defaults
	 */

	if len(client.ResponseTypes) == 0 {
		client.ResponseTypes = []OidcClientResponseType{
			{
				ResponseType: "code",
			},
		}
	}

	if len(client.GrantTypes) == 0 {
		client.GrantTypes = []OidcClientGrantType{
			{
				GrantType: "authorization_code",
			},
		}
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

func (client *OidcClient) GetRedirectURIs() []string {
	res := []string{}
	for _, redirectURI := range client.RedirectURIs {
		res = append(res, redirectURI.RedirectURI)
	}

	return res
}

func (client *OidcClient) GetRequestURIs() []string {
	res := []string{}
	for _, requestURI := range client.RequestURIs {
		res = append(res, requestURI.RequestURI)
	}

	return res
}

func (client *OidcClient) GetResponseTypes() []string {
	res := []string{}
	for _, responseType := range client.ResponseTypes {
		res = append(res, responseType.ResponseType)
	}

	return res
}

func (client *OidcClient) GetGrantTypes() []string {
	res := []string{}
	for _, grantType := range client.GrantTypes {
		res = append(res, grantType.GrantType)
	}

	return res
}

func (client *OidcClient) CanRedirectTo(uri string) bool {
	return utils.Contains(client.GetRedirectURIs(), uri)
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

	res := DB.First(&client, "id = ?", clientId)
	if res.Error != nil {
		log.Panicf("Error while Loading client: %v", res.Error)
	}

	return client
}

func (client *OidcClient) Validate() error {
	redirectURIs := client.GetRedirectURIs()
	if err := validateURIArray(redirectURIs); err != nil {
		return err
	}

	grantTypes := client.GetGrantTypes()

	responseGrantsDependencies := map[string][]string{
		"code":                {"authorization_code"},
		"id_token":            {"implicit"},
		"token id_token":      {"implicit"},
		"code id_token":       {"authorization_code", "implicit"},
		"code token":          {"authorization_code", "implicit"},
		"code token id_token": {"authorization_code", "implicit"},
	}

	for _, responseType := range client.GetResponseTypes() {
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

	if client.JWKsURI != "" {
		if err := validateURI(client.JWKsURI); err != nil {
			return err
		}

		if len(client.JWKs) != 0 {
			return fmt.Errorf("RP cannot provide both jwks_uri and jwks")
		}
	}

	if client.InitiateLoginURI != "" {
		if err := validateURI(client.InitiateLoginURI); err != nil {
			return err
		}
	}

	if err := validateURIArray(client.GetRequestURIs()); err != nil {
		return err
	}

	return nil
}
