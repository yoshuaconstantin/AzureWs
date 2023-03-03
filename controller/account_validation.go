package controller

import (
	"net/http"

	Auth "AzureWS/globalvariable/authenticator"
	"AzureWS/validation"
)

func ValidatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")
	passwordParam := queryParams.Get("password")

	GetUserIdAuth, AunthStatus, errAunth := Auth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	storedPassword, errStrPswd := validation.ValidateGetStoredPasswordByUserId(GetUserIdAuth)

	if errStrPswd != nil {
		http.Error(w, errStrPswd.Error(), http.StatusUnauthorized)
		return
	}

	PasswordValidation, errPassvalidate := validation.ValidateUserPassword(passwordParam, storedPassword)

	if errPassvalidate != nil {
		http.Error(w, errPassvalidate.Error(), http.StatusBadRequest)
		return
	}

	if !PasswordValidation {

		http.Error(w, string("Password not match!"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
