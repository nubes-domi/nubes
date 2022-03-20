package oidc

import (
	"net/http"
	"nubes/sum/utils"

	"github.com/gin-gonic/gin"
)

type OpenidConfiguration struct {
	Issuer                                     string   `json:"issuer"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                              string   `json:"token_endpoint"`
	UserinfoEndpoint                           string   `json:"userinfo_endpoint,omitempty"`
	JwksURI                                    string   `json:"jwks_uri"`
	RegistrationEndpoint                       string   `json:"registration_endpoint,omitempty"`
	ScopesSupported                            []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported                     []string `json:"response_types_supported"`
	ResponseModesSupported                     []string `json:"response_modes_supported,omitempty"`
	GrantTypesSupported                        []string `json:"grant_types_supported,omitempty"`
	ACRValuesSupported                         []string `json:"acr_values_supported,omitempty"`
	SubjectTypesSupported                      []string `json:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported           []string `json:"id_token_signing_alg_values_supported"`
	IdTokenEncryptionAlgValuesSupported        []string `json:"id_token_encryption_alg_values_supported,omitempty"`
	IdTokenEncryptionEncValuesSupported        []string `json:"id_token_encryption_enc_values_supported,omitempty"`
	UserinfoSigningAlgValuesSupported          []string `json:"userinfo_signing_alg_values_supported,omitempty"`
	UserinfoEncryptionAlgValuesSupported       []string `json:"userinfo_encryption_alg_values_supported,omitempty"`
	UserinfoEncryptionEncValuesSupported       []string `json:"userinfo_encryption_enc_values_supported,omitempty"`
	RequestObjectSigningAlgValuesSupported     []string `json:"request_object_signing_alg_values_supported,omitempty"`
	RequestObjectEncryptionAlgValuesSupported  []string `json:"request_object_encryption_alg_values_supported,omitempty"`
	RequestObjectEncryptionEncValuesSupported  []string `json:"request_object_encryption_enc_values_supported,omitempty"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	DisplayValuesSupported                     []string `json:"display_values_supported,omitempty"`
	ClaimTypesSupported                        []string `json:"claim_types_supported,omitempty"`
	ClaimsSupported                            []string `json:"claims_supported,omitempty"`
	ServiceDocumentation                       string   `json:"service_documentation,omitempty"`
	ClaimsLocalesSupported                     []string `json:"claims_locales_supported,omitempty"`
	UILocalesSupported                         []string `json:"ui_locales_supported,omitempty"`
	ClaimsParameterSupported                   bool     `json:"claims_parameter_supported,omitempty"`
	RequestParameterSupported                  bool     `json:"request_parameter_supported,omitempty"`
	RequestURIParamterSupported                bool     `json:"request_uri_parameter_supported,omitempty"`
	RequireRequestURIRegistration              bool     `json:"require_request_uri_registration,omitempty"`
	OpPolicyURI                                string   `json:"op_policy_uri,omitempty"`
	OpTosURI                                   string   `json:"op_tos_uri,omitempty"`
}

func Discovery(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, OpenidConfiguration{
		Issuer:                                 baseURI(c),
		AuthorizationEndpoint:                  baseURI(c) + "/openid/authorization",
		TokenEndpoint:                          baseURI(c) + "/openid/token",
		UserinfoEndpoint:                       baseURI(c) + "/openid/userinfo",
		JwksURI:                                baseURI(c) + "/openid/jwks",
		RegistrationEndpoint:                   baseURI(c) + "/openid/registration",
		ScopesSupported:                        []string{"openid", "profile", "email", "address", "phone"},
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

func Jwks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, utils.JWKPublicSet)
}

func baseURI(c *gin.Context) string {
	return "https://" + c.Request.Host
}