package main

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
)

func baseURI(c *gin.Context) string {
	return "https://" + c.Request.Host
}

func openidConfiguration(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, OpenidConfiguration{
		Issuer:                           baseURI(c),
		AuthorizationEndpoint:            baseURI(c) + "/openid/authorization",
		TokenEndpoint:                    baseURI(c) + "/openid/token",
		UserinfoEndpoint:                 baseURI(c) + "/openid/userinfo",
		JWKsURI:                          baseURI(c) + "/openid/jwks",
		RegistrationEndpoint:             baseURI(c) + "/openid/registration",
		ScopesSupported:                  []string{"openid"},
		ResponseTypesSupported:           []string{"code", "id_token", "token id_token"},
		ResponseModesSupported:           []string{"query", "fragment"},
		GrantTypesSupported:              []string{"authorization_code", "implicit"},
		ACRValuesSupported:               []string{},
		SubjectTypesSupported:            []string{"public"},
		IdTokenSigningAlgValuesSupported: []string{"RS256"},
	})
}

var JWKSet jwk.Set

func jwks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, JWKSet)
}

type RegistrationResponse struct {
	ClientID                string `json:"client_id"`
	ClientSecret            string `json:"client_secret,omitempty"`
	RegistrationAccessToken string `json:"registration_access_token,omitempty"`
	RegistrationClientURI   string `json:"registration_client_uri,omitempty"`
	ClientIDIssuedAt        string `json:"client_id_issued_at,omitempty"`
	ClientSecretExpiresAt   int    `json:"client_secret_expires_at"`
}

func registration(c *gin.Context) {
	c.IndentedJSON(http.StatusCreated, RegistrationResponse{
		ClientID:              "my-client-id",
		ClientSecret:          "a-fake-secret",
		ClientSecretExpiresAt: 0,
	})
}

func authorizationStart(c *gin.Context) {
	c.HTML(http.StatusOK, "new_session.html", gin.H{
		"scope":         c.Query("scope"),
		"response_type": c.Query("response_type"),
		"client_id":     c.Query("client_id"),
		"redirect_uri":  c.Query("redirect_uri"),
		"state":         c.Query("state"),
		"response_mode": c.Query("response_mode"),
		"nonce":         c.Query("nonce"),
		"display":       c.Query("display"),
		"prompt":        c.Query("prompt"),
		"max_age":       c.Query("max_age"),
		"ui_locales":    c.Query("ui_locales"),
		"id_token_hint": c.Query("id_token_hint"),
		"login_hint":    c.Query("login_hint"),
		"acr_values":    c.Query("acr_values"),
	})
}

func random256() string {
	buffer := make([]byte, 24)
	rand.Read(buffer)

	return base64.URLEncoding.EncodeToString(buffer)
}

func authorizationSubmit(c *gin.Context) {
	redirectURI := c.PostForm("redirect_uri")
	state := c.PostForm("state")

	if !strings.Contains(redirectURI, "?") {
		redirectURI += "?"
	} else {
		redirectURI += "&"
	}

	redirectURI += "state=" + state
	redirectURI += "&code=" + random256() + ";" + c.PostForm("nonce")

	c.Redirect(http.StatusFound, redirectURI)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	IdToken      string `json:"id_token"`
}

type IdToken struct {
	Iss      string
	Sub      string
	Aud      string
	Exp      int
	Iat      int
	AuthTime int
	Nonce    string
	Acr      []string
	Amr      []string
	Azp      []string
}

func token(c *gin.Context) {
	id := openid.New()
	id.Set(jwt.IssuerKey, baseURI(c))
	id.Set(jwt.SubjectKey, "me")
	id.Set(jwt.AudienceKey, "my-client-id")
	id.Set(jwt.ExpirationKey, 1647423942)
	id.Set(jwt.IssuedAtKey, 123123)

	nonce := strings.Split(c.PostForm("code"), ";")[1]
	id.Set("nonce", nonce)

	signed, err := jwt.Sign(id, jwa.RS256, RSAKey)
	if err != nil {
		log.Panicf("Could not sign IDToken: %v", err)
	}

	c.IndentedJSON(http.StatusCreated, TokenResponse{
		AccessToken:  random256(),
		TokenType:    "Bearer",
		RefreshToken: "",
		ExpiresIn:    3600,
		IdToken:      string(signed[:]),
	})
}

type UserinfoResponse struct {
}

func userinfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, UserinfoResponse{})
}
