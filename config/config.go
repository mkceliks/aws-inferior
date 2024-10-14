package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

var (
	ClientID     = "734r4pj66slej1sl7rqmgsmtqi"                          // Cognito Client ID
	ClientSecret = "ho37damcm3esla700re2k88t23d95797qv82jgn26gd7enavr1b" // Cognito Client Secret
)

// SECRET_HASH hesaplama fonksiyonu
func CalculateSecretHash(clientSecret, clientID, username string) string {
	key := []byte(clientSecret)
	message := []byte(username + clientID)

	hmacHash := hmac.New(sha256.New, key)
	hmacHash.Write(message)

	secretHash := base64.StdEncoding.EncodeToString(hmacHash.Sum(nil))
	return secretHash
}
