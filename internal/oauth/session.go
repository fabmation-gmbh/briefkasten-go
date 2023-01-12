package oauth

// UserSessionTokenClaims represents the ID Token claim object.
type UserSessionTokenClaims struct {
	// Exp int `json:"exp"`
	// Iat               int    `json:"iat"`
	// AuthTime          int    `json:"auth_time"`
	// Jti               string `json:"jti"`
	// Iss               string `json:"iss"`
	// Aud               string `json:"aud"`
	// Sub               string `json:"sub"`
	// Typ               string `json:"typ"`
	// Azp               string `json:"azp"`
	// Nonce             string `json:"nonce"`
	// SessionState      string `json:"session_state"`
	// AtHash            string `json:"at_hash"`
	// Acr               string `json:"acr"`
	// Sid               string `json:"sid"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}
