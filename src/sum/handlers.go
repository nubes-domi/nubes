package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"nubes/sum/db"
	"nubes/sum/utils"
	"strconv"
	"strings"
	"time"

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
		Issuer:                                 baseURI(c),
		AuthorizationEndpoint:                  baseURI(c) + "/openid/authorization",
		TokenEndpoint:                          baseURI(c) + "/openid/token",
		UserinfoEndpoint:                       baseURI(c) + "/openid/userinfo",
		JwksURI:                                baseURI(c) + "/openid/jwks",
		RegistrationEndpoint:                   baseURI(c) + "/openid/registration",
		ScopesSupported:                        []string{"openid"},
		ResponseTypesSupported:                 []string{"code", "id_token", "token id_token"},
		ResponseModesSupported:                 []string{"query", "fragment"},
		GrantTypesSupported:                    []string{"authorization_code", "implicit"},
		ACRValuesSupported:                     []string{},
		SubjectTypesSupported:                  []string{"public", "pairwise"},
		IdTokenSigningAlgValuesSupported:       []string{"none", "RS256"},
		UserinfoSigningAlgValuesSupported:      []string{"none", "RS256"},
		TokenEndpointAuthMethodsSupported:      []string{"none", "client_secret_post", "client_secret_basic", "client_secret_jwt", "private_key_jwt"},
		RequestParameterSupported:              true,
		RequestURIParamterSupported:            true,
		RequestObjectSigningAlgValuesSupported: []string{"none", "RS256"},
		ClaimsSupported:                        []string{"sub", "iss", "auth_time", "name", "email", "locale", "zoneinfo", "preferred_username", "profile", "picture", "phone_number"},
		ClaimsParameterSupported:               true,
	})
}

var JWKSet jwk.Set
var JWKPublicSet jwk.Set

func jwks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, JWKPublicSet)
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
	client := db.BuildOpenIDClient(c)
	db.DB.OidcClients().Create(&client)

	c.IndentedJSON(http.StatusCreated, RegistrationResponse{
		ClientID:              client.ID,
		ClientSecret:          "a-fake-secret",
		ClientSecretExpiresAt: 0,
	})
}

func buildRedirection(baseURI string, additionalParams map[string]string) string {
	if len(additionalParams) == 0 {
		return baseURI
	}

	if !strings.Contains(baseURI, "?") {
		baseURI += "?"
	} else {
		baseURI += "&"
	}

	for k, v := range additionalParams {
		baseURI += url.QueryEscape(k) + "=" + url.QueryEscape(v) + "&"
	}

	return baseURI
}

type AuthorizationRequest struct {
	ResponseType string
	RedirectURI  string
	State        string
	Scope        string
	ClientID     string
	ResponseMode string
	Nonce        string
	Display      string
	Prompt       string
	MaxAge       int
	UILocales    string
	IDTokenHint  string
	LoginHint    string
	ACRValues    string
	Claims       map[string]interface{}
}

func buildAuthorizationRequest(c *gin.Context) AuthorizationRequest {
	maxAgeStr := c.Query("max_age")
	maxAge, err := strconv.Atoi(maxAgeStr)
	if maxAgeStr != "" && err != nil {
		log.Panicf("Could not parse max_age: %s", maxAgeStr)
	}

	authRequest := AuthorizationRequest{
		ResponseType: c.Query("response_type"),
		RedirectURI:  c.Query("redirect_uri"),
		State:        c.Query("state"),
		Scope:        c.Query("scope"),
		ClientID:     c.Query("client_id"),
		ResponseMode: c.Query("response_mode"),
		Nonce:        c.Query("nonce"),
		Display:      c.Query("display"),
		Prompt:       c.Query("prompt"),
		MaxAge:       maxAge,
		UILocales:    c.Query("ui_locales"),
		IDTokenHint:  c.Query("id_token_hint"),
		LoginHint:    c.Query("login_hint"),
		ACRValues:    c.Query("acr_values"),
	}

	var requestObject []byte
	if c.Query("request_uri") != "" {
		if c.Query("request") != "" {
			log.Panicf("Cannot use both request_uri and request")
		}

		resp, err := http.Get(c.Query("request_uri"))
		if err != nil {
			log.Panicf("%v", err)
		}

		requestObject, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicf("%v", err)
		}
	}

	if c.Query("request") != "" {
		requestObject = []byte(c.Query("request"))
	}

	if len(requestObject) > 0 {
		token, err := jwt.Parse(requestObject) //, jwt.WithVerify(jwa.RS256, RSAKey.PublicKey))
		if err != nil {
			log.Panicf("%v", err)
		}

		if val, ok := token.Get("response_type"); ok {
			authRequest.ResponseType = val.(string)
		}

		if val, ok := token.Get("redirect_uri"); ok {
			authRequest.RedirectURI = val.(string)
		}

		if val, ok := token.Get("state"); ok {
			authRequest.State = val.(string)
		}

		if val, ok := token.Get("scope"); ok {
			authRequest.Scope = val.(string)
		}

		if val, ok := token.Get("client_id"); ok {
			authRequest.ClientID = val.(string)
		}

		if val, ok := token.Get("response_mode"); ok {
			authRequest.ResponseMode = val.(string)
		}

		if val, ok := token.Get("nonce"); ok {
			authRequest.Nonce = val.(string)
		}

		if val, ok := token.Get("display"); ok {
			authRequest.Display = val.(string)
		}

		if val, ok := token.Get("prompt"); ok {
			authRequest.Prompt = val.(string)
		}

		if val, ok := token.Get("max_age"); ok {
			authRequest.MaxAge = val.(int)
		}

		if val, ok := token.Get("ui_locales"); ok {
			authRequest.UILocales = val.(string)
		}

		if val, ok := token.Get("id_token_hint"); ok {
			authRequest.IDTokenHint = val.(string)
		}

		if val, ok := token.Get("login_hint"); ok {
			authRequest.LoginHint = val.(string)
		}

		if val, ok := token.Get("acr_values"); ok {
			authRequest.ACRValues = val.(string)
		}
	}

	return authRequest
}

func authorizationStart(c *gin.Context) {
	request := buildAuthorizationRequest(c)

	client := db.LoadClient(request.ClientID)
	if !utils.Contains(client.RedirectURIs, request.RedirectURI) {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"error": "The application specified a Redirect URI which has not been whitelisted.",
		})
		return
	}

	if request.ResponseType == "" {
		c.Redirect(302, buildRedirection(request.RedirectURI, map[string]string{
			"error": "invalid_request",
			"state": request.State,
		}))
	} else {
		logo := client.GetLogoURI("")
		policyUri := client.GetPolicyURI("")
		tosUri := client.GetTosURI("")
		c.HTML(http.StatusOK, "new_session.html", gin.H{
			"client":  client,
			"request": request,
			"logo":    logo,
			"policy":  policyUri,
			"tos":     tosUri,
		})
	}
}

func random128() string {
	return base64.URLEncoding.EncodeToString(utils.RandBytes(12))
}

func generateAuthorizationCode(subject, nonce, clientId string) string {
	token, err := jwt.NewBuilder().
		Expiration(time.Now().Add(time.Minute*10)).
		JwtID(random128()).
		Subject(subject).
		Audience([]string{"token"}).
		Issuer("authorization").
		Claim("nonce", nonce).
		Claim("client_id", clientId).
		Build()

	if err != nil {
		log.Panicf("Could not build a JWT: %v", err)
	}

	str, err := jwt.Sign(token, jwa.RS256, RSAKey)
	if err != nil {
		log.Panicf("Could not build a JWT: %v", err)
	}

	return string(str[:])
}

func authorizationSubmit(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	_, err := db.DB.Users().FindByCredentials(username, password)
	if err != nil {
		log.Panicf("Could not find the user")
	}

	redirectURI := c.PostForm("redirect_uri")
	state := c.PostForm("state")

	if !strings.Contains(redirectURI, "?") {
		redirectURI += "?"
	} else {
		redirectURI += "&"
	}

	c.Redirect(http.StatusFound, buildRedirection(redirectURI, map[string]string{
		"state": state,
		"code":  generateAuthorizationCode("123", c.PostForm("nonce"), c.PostForm("client_id")),
	}))
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	IdToken      string `json:"id_token"`
}

func token(c *gin.Context) {
	id := openid.New()
	id.Set(jwt.IssuerKey, baseURI(c))
	id.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	id.Set(jwt.IssuedAtKey, time.Now())

	token, err := jwt.Parse([]byte(c.PostForm("code")), jwt.WithVerify(jwa.RS256, RSAKey.PublicKey))
	if err != nil {
		log.Panicf("Invalid authorization code: %v", err)
	}

	id.Set(jwt.SubjectKey, token.Subject())
	clientId, _ := token.Get("client_id")
	id.Set(jwt.AudienceKey, clientId)
	nonce, _ := token.Get("nonce")
	id.Set("nonce", nonce)

	client := db.LoadClient(clientId.(string))
	signed := utils.JwtSign(id, client.IDTokenSignedResponseAlg, JWKSet)

	c.IndentedJSON(http.StatusCreated, TokenResponse{
		AccessToken:  random128(),
		TokenType:    "Bearer",
		RefreshToken: "",
		ExpiresIn:    3600,
		IdToken:      string(signed[:]),
	})
}

type UserinfoResponse struct {
	Sub string `json:"sub"`
}

func userinfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, UserinfoResponse{
		Sub: "123",
	})
}
