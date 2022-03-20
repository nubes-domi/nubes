package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

var JWKSet jwk.Set
var JWKPublicSet jwk.Set

func getKeyForSigning(algorithm string) jwk.Key {
	for it := JWKSet.Iterate(context.Background()); it.Next(context.Background()); {
		pair := it.Pair()
		key := pair.Value.(jwk.Key)

		switch algorithm {
		case "RS256", "RS384", "RS512", "PS256", "PS384", "PS512":
			if key.KeyType() == "RSA" {
				return key
			}
		case "HS256", "HS384", "HS512":
			if key.KeyType() == "oct" {
				return key
			}
		case "ES256", "ES384", "ES512":
			if key.KeyType() == "EC" {
				return key
			}
		}
	}

	log.Panicf("Could not find appropriate key for signing algorithm %s", algorithm)
	return nil
}

func JwtSign(token jwt.Token, algorithm string) string {
	switch algorithm {
	case "none":
		marshalled, err := json.Marshal(token)
		if err != nil {
			log.Panicf("Could not marshal token to JSON")
		}

		encoded := base64.RawURLEncoding.EncodeToString(marshalled)

		// {"alg":"none"}
		return "eyJhbGciOiJub25lIn0." + encoded + "."
	case "RS256", "RS384", "RS512", "HS256", "HS384", "HS512", "ES256", "ES384", "ES512", "PS256", "PS384", "PS512":
		key := getKeyForSigning(algorithm)

		jwsHeaders := jws.NewHeaders()
		jwsHeaders.Set("kid", key.KeyID())

		signed, err := jwt.Sign(token, jwa.SignatureAlgorithm(algorithm), key, jwt.WithJwsHeaders(jwsHeaders))
		if err != nil {
			log.Panicf("Could not sign JWT: %v", err)
		}

		return string(signed)
	default:
		log.Panicf("Unrecognized signing algorithm %s", algorithm)
		return ""
	}
}

func JwtVerify(serialized string) (jwt.Token, error) {
	return jwt.Parse([]byte(serialized), jwt.WithKeySet(JWKPublicSet), jwt.InferAlgorithmFromKey(true))
}

func PrepareKeys() {
	JWKSet = jwk.NewSet()
	JWKPublicSet = jwk.NewSet()

	rsaPriv, rsaPub := RSAKeypair("sig")
	ecPriv, ecPub := ECKeypair("sig")
	okpPriv, okpPub := OKPKeypair("sig")

	JWKSet.Add(rsaPriv)
	JWKSet.Add(ecPriv)
	JWKSet.Add(okpPriv)

	JWKPublicSet.Add(rsaPub)
	JWKPublicSet.Add(ecPub)
	JWKPublicSet.Add(okpPub)
}

func RSAKeypair(use string) (jwk.Key, jwk.Key) {
	rsa, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicf("failed to generate new RSA key: %v\n", err)
	}

	private, err := jwk.New(rsa)
	if err != nil {
		log.Panicf("failed to wrap RSA key into JWK: %v\n", err)
	}

	private.Set("use", use)
	private.Set("kid", Sha256String(rsa.N.Bytes())[:8])

	public, err := private.PublicKey()
	if err != nil {
		log.Panicf("failed to get public key: %v", err)
	}

	return private, public
}

func ECKeypair(use string) (jwk.Key, jwk.Key) {
	ec, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panicf("failed to generate new EC key: %v\n", err)
	}

	private, err := jwk.New(ec)
	if err != nil {
		log.Panicf("failed to wrap EC key into JWK: %v\n", err)
	}

	private.Set("use", use)
	private.Set("kid", Sha256String(append(ec.X.Bytes(), ec.Y.Bytes()...))[:8])

	public, err := private.PublicKey()
	if err != nil {
		log.Panicf("failed to get public key: %v", err)
	}

	return private, public
}

func OKPKeypair(use string) (jwk.Key, jwk.Key) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Panicf("failed to generate new EC key: %v\n", err)
	}

	private, err := jwk.New(priv)
	if err != nil {
		log.Panicf("failed to wrap EC key into JWK: %v\n", err)
	}

	public, err := jwk.New(pub)
	if err != nil {
		log.Panicf("failed to wrap EC key into JWK: %v\n", err)
	}

	public.Set("use", use)
	private.Set("use", use)
	public.Set("kid", Sha256String(pub)[:8])
	private.Set("kid", Sha256String(pub)[:8])

	return private, public
}
