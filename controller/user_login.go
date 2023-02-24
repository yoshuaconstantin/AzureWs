package controller

import (
	jwttoken "AzureWS/JWTTOKEN"
	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/module"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"AzureWS/session"
	"AzureWS/validation"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux" 
	_ "github.com/lib/pq"    
)

var ipMap = make(map[string]int)

var blockedIPs = make(map[string]time.Time)

func isBlocked(ipAddr string) bool {
	blockedAt, found := blockedIPs[ipAddr]
	if found && time.Since(blockedAt) < 4*time.Hour {
		return true
	}
	return false
}

// Block an IP address
func blockIP(ipAddr string) {
	blockedIPs[ipAddr] = time.Now()
}

// Reset a blocked IP address after 4 hours
func resetIP(ipAddr string) {
	delete(blockedIPs, ipAddr)
}

func CreateNewAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var user models.UserModel

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("\nCREATE NEW USER - Cannot get request body : %v\n", err)
	}

	GetTokenAndJwt, errInsert := module.AddUserToDB(user)

	if errInsert != nil {
		http.Error(w, errInsert.Error(), http.StatusInternalServerError)
		return
	}

	var response response.ResponseUserLoginWithJWT
	response.Status = http.StatusOK
	response.Message = "Data user baru telah di tambahkan"
	response.JwtToken = GetTokenAndJwt.JWT
	response.Token = GetTokenAndJwt.Token

	// kirim response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// AmbilUser mengambil single data dengan parameter id
func GetSnglUsr(w http.ResponseWriter, r *http.Request) {
	// kita set headernya
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// dapatkan idUser dari parameter request, keynya adalah "id"
	params := mux.Vars(r)

	// konversi id dari tring ke int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Tidak bisa mengubah dari string ke int.  %v", err)
	}

	// memanggil models GetSingleUser dengan parameter id yg nantinya akan mengambil single data
	user, err := module.GetSingleUser(int64(id))

	if err != nil {
		log.Fatalf("Tidak bisa mengambil data User. %v", err)
	}

	// kirim response
	json.NewEncoder(w).Encode(user)
}

// Login user
func LoginAccount(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ipAddr := r.RemoteAddr

	IsIpBlocked := isBlocked(ipAddr)

	if IsIpBlocked {
		http.Error(w, "Your IP has been blocked try again in 4 more hours.", http.StatusUnauthorized)
		return
	} else {
		resetIP(ipAddr)
	}

	if ipMap[ipAddr] >= 5 {
		// Return an error response indicating that the IP address has been blocked
		blockIP(ipAddr)

		http.Error(w, "Too many failed login attempts. Your IP has been blocked.", http.StatusUnauthorized)
		return
	}

	var loginData request.RequestLoginData
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := loginData.Username
	password := loginData.Password

	storedPassword, errStrdPswd := validation.ValidateGetStoredPassword(username)

	if errStrdPswd != nil {
		http.Error(w, errStrdPswd.Error(), http.StatusInternalServerError)
		return
	}

	PasswordValidation, errPassvalidate := validation.ValidateUserPassword(password, storedPassword)

	fmt.Printf("Password Validation = %v\n", PasswordValidation)
	if errPassvalidate != nil {
		// kirim respon 400 kalau ada error
		ipMap[ipAddr]++

		attemptsLeft := 5 - ipMap[ipAddr]

		errorMsg := fmt.Sprintf("Password not match. You have %d attempts left", attemptsLeft)

		http.Error(w, errorMsg, http.StatusConflict)
		return
	}

	if PasswordValidation {

		token, err := validation.Validate(username, storedPassword)

		if err == nil {
			if token != "" {

				GetUserID, errGetUuid := validation.ValidateTokenGetUuid(token)

				if errGetUuid != nil {
					http.Error(w, errGetUuid.Error(), http.StatusUnauthorized)
					return
				}

				// Generate JWT Token after Succesfully passed the aunth system
				GenerateJwtToken, errGenerate := jwttoken.GenerateToken(GetUserID)

				if errGenerate != nil {
					http.Error(w, errGenerate.Error(), http.StatusInternalServerError)
					return
				}

				ReNewSession, errAddSession := session.ReNewSessionLogin(GetUserID)

				if errAddSession != nil {
					http.Error(w, errAddSession.Error(), http.StatusInternalServerError)
					return
				}

				if !ReNewSession {
					http.Error(w, "Failed to refresh session", http.StatusInternalServerError)
					return
				}

				fmt.Printf("\nPASSWORD VALIDATION - User Token = %v\n", token)

				ipMap[ipAddr] = 0
				// Send Token and JWT token
				var response response.ResponseUserTokenAndJwt
				response.Status = http.StatusOK
				response.Message = "Success"
				response.Token = token
				response.JwtToken = GenerateJwtToken

				w.WriteHeader(http.StatusOK)

				json.NewEncoder(w).Encode(response)
				return

			} else {

				// kirim respon 500 kalau token kosong
				var response response.GeneralResponseNoData
				response.Status = http.StatusInternalServerError
				response.Message = "No token"

				w.WriteHeader(http.StatusInternalServerError)

				json.NewEncoder(w).Encode(response)
				return

			}
		} else {

			var response response.GeneralResponseNoData
			response.Status = http.StatusBadRequest
			response.Message = "User not found "

			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(response)
			return
		}
	} else {

		ipMap[ipAddr]++

		attemptsLeft := 5 - ipMap[ipAddr]

		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("Password validating failed. You have %d attempts left.\n", attemptsLeft)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func GetAllUsr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Call the GetAllUser method from the models package
	users, err := module.GetAllUser()

	if err != nil {
		log.Fatalf("Unable to retrieve data. %v", err)

		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving data"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	var response response.ResponseUserLogin
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = users

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update User Password func
func UpdateAccountPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var updatePswdModel request.RequestChangePasswordData

	err := json.NewDecoder(r.Body).Decode(&updatePswdModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, updatePswdModel.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	_, errUpdatePswd := module.UpdatePasswordUserFromDB(GetUserIdAunth, updatePswdModel.Password)

	if errUpdatePswd != nil {
		http.Error(w, errUpdatePswd.Error(), http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Message = "Succes"
	response.Status = http.StatusOK

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete User func
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var tokenModel request.RequestTokenData

	err := json.NewDecoder(r.Body).Decode(&tokenModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenModel.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	_, errDeleteUser := module.RemoveUserDB(GetUserIdAunth)

	if errDeleteUser != nil {
		http.Error(w, errDeleteUser.Error(), http.StatusInternalServerError)
		return

	}

	var response response.GeneralResponseNoData
	response.Message = "Delete user operation success"
	response.Status = http.StatusOK

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logout User func
func LogoutAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var tokenModel request.RequestTokenData

	err := json.NewDecoder(r.Body).Decode(&tokenModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenModel.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	logoutUser, errLogout := module.LogoutUser(GetUserIdAunth)

	if errLogout != nil {
		http.Error(w, errLogout.Error(), http.StatusInternalServerError)
		return
	}

	if !logoutUser {
		http.Error(w, "Error when trying to logout", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Message = "Logout Succes"
	response.Status = http.StatusOK

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Refresh JWT, place this on init
func RefrshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var tokenModel request.RequestTokenData

	err := json.NewDecoder(r.Body).Decode(&tokenModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenModel.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	// Now check and refresh the token
	RefreshJWT, errRefreshJWT := jwttoken.RefreshToken(tokenModel.Token, GetUserIdAunth)

	if errRefreshJWT != nil {
		http.Error(w, errRefreshJWT.Error(), http.StatusInternalServerError)
		return
	}

	var response response.ResponseUserLoginWithJWT
	response.Message = "Token Refreshed"
	response.Status = http.StatusOK
	response.JwtToken = RefreshJWT

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Testing Generate JWT
func TestGenerateJwt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var loginData request.RequestLoginData
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GenerateJwt, errGen := jwttoken.GenerateToken(loginData.Username)

	if errGen != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response response.ResponseUserLoginWithJWT
	response.Status = http.StatusOK
	response.Message = "Testing Generate JWT"
	response.JwtToken = GenerateJwt

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Testing Verify JWT
func TestVerifyJwt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// If the authorization header is empty, return an error
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Missing authorization header"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(tokenString)

	if erroCheckJWt != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = erroCheckJWt.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckJwtTokenValidation {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Unauthorized user"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
