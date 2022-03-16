package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
)

var RSAKey *rsa.PrivateKey

func PrepareKeys() {
	JWKSet = jwk.NewSet()

	var err error
	RSAKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicf("failed to generate new RSA privatre key: %s\n", err)
	}

	key, err := jwk.New(RSAKey)
	if err != nil {
		log.Panicf("failed to create symmetric key: %s\n", err)
	}

	publicKey, err := key.PublicKey()
	publicKey.Set("use", "sig")
	publicKey.Set("kid", "deadbeef")
	if err != nil {
		log.Panicf("expected jwk.SymmetricKey, got %T\n", key)
		return
	}

	JWKSet.Add(publicKey)
}

func main() {
	PrepareKeys()

	router := gin.Default()
	router.LoadHTMLFiles("new_session.html")

	router.GET("/.well-known/openid-configuration", openidConfiguration)
	router.GET("/openid/jwks", jwks)
	router.POST("/openid/registration", registration)

	router.GET("/openid/authorization", authorizationStart)
	router.POST("/openid/authorization", authorizationSubmit)

	router.POST("/openid/token", token)

	router.GET("/openid/userinfo", userinfo)

	router.Run("localhost:8080")
}
