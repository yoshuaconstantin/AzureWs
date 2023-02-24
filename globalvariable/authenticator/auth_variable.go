package authenticator

import (
	jwttoken "AzureWS/JWTTOKEN"
	"AzureWS/session"
	"AzureWS/validation"
	"fmt"
	"net/http"
	"strings"
)

// Authenticatior main function to securely authenticate request also shorten usage code
func SecureAuthenticator(w http.ResponseWriter, r *http.Request, Token string) (string, int, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", http.StatusBadRequest, fmt.Errorf("%s", "Missing authorization header")
	}

	userId, errGetUuid := validation.ValidateTokenGetUuid(Token)

	if errGetUuid != nil {
		return "", http.StatusInternalServerError, errGetUuid
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(tokenString)

	if erroCheckJWt != nil {
		return "", http.StatusUnauthorized, erroCheckJWt
	}

	if !CheckJwtTokenValidation {
		return "", http.StatusUnauthorized, fmt.Errorf("%s", "Unauthorized user")
	}

	SessionValidation, errSessionCheck := session.CheckSessionInside(userId)

	if errSessionCheck != nil {
		return "", http.StatusForbidden, errSessionCheck
	}

	if !SessionValidation {
		return "", http.StatusUnauthorized, fmt.Errorf("%s", "Session Expired")
	}

	return userId, 0, nil
}
