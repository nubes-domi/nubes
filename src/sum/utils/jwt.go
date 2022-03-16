package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

func getKeyForSigning(set jwk.Set, algorithm string) jwk.Key {
	for it := set.Iterate(context.Background()); it.Next(context.Background()); {
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

func JwtSign(token jwt.Token, algorithm string, keyBundle jwk.Set) string {
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
		key := getKeyForSigning(keyBundle, algorithm)

		jwsHeaders := jws.NewHeaders()
		jwsHeaders.Set("kid", key.KeyID())

		signed, err := jwt.Sign(token, jwa.SignatureAlgorithm(algorithm), key, jwt.WithJwsHeaders(jwsHeaders))
		if err != nil {
			log.Panicf("Could not sign IDToken: %v", err)
		}

		return string(signed)
	default:
		log.Panicf("Unrecognized signing algorithm %s", algorithm)
		return ""
	}
}
