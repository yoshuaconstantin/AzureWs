package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	jwttoken "AzureWS/JWTTOKEN"
	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/globalvariable/constant"
	"AzureWS/logging"
	"AzureWS/module"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"AzureWS/session"
	"AzureWS/validation"

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

	GetTokenAndJwt, errInsert := module.CreateAccountToDB(user)

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

	storedPassword, errStrdPswd := validation.ValidateGetStoredPasswordByUsername(username)

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

	if !PasswordValidation {
		ipMap[ipAddr]++

		attemptsLeft := 5 - ipMap[ipAddr]

		errorMsg := fmt.Sprintf("Password not match. You have %d attempts left", attemptsLeft)

		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	token, errToken := validation.ValidateGenerateNewToken(username, password)

	if errToken != nil {
		http.Error(w, errToken.Error(), http.StatusUnauthorized)
	}

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

	ipMap[ipAddr] = 0

	var response response.ResponseUserTokenAndJwt
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Token = token
	response.JwtToken = GenerateJwtToken

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
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

		logging.InsertLog(r, constant.UserUpdatePassword, err.Error(), "", http.StatusBadRequest, 3, 2)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, updatePswdModel.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.UserUpdatePassword, errAunth.Error(), updatePswdModel.Token, AunthStatus, 3, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	UpdatePswd, errUpdatePswd := module.UpdatePasswordAccountFromDB(GetUserIdAunth, updatePswdModel.Password)

	if errUpdatePswd != nil {

		logging.InsertLog(r, constant.UserUpdatePassword, errUpdatePswd.Error(), updatePswdModel.Token, http.StatusInternalServerError, 3, 3)

		http.Error(w, errUpdatePswd.Error(), http.StatusInternalServerError)
		return
	}

	if !UpdatePswd {

		logging.InsertLog(r, constant.UserUpdatePassword, "Cannot change password", updatePswdModel.Token, http.StatusInternalServerError, 3, 3)

		http.Error(w, "Cannot change password", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Message = "Succes"
	response.Status = http.StatusOK

	logging.InsertLog(r, constant.UserUpdatePassword, "", updatePswdModel.Token, http.StatusOK, 3, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete User func
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {

		logging.InsertLog(r, constant.UserDeleteAccount, errAunth.Error(), tokenParam, AunthStatus, 4, 3)
		
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	_, errDeleteUser := module.RemoveAccountFromDB(GetUserIdAunth)

	if errDeleteUser != nil {

		logging.InsertLog(r, constant.UserDeleteAccount, errDeleteUser.Error(), tokenParam, http.StatusInternalServerError, 4, 3)

		http.Error(w, errDeleteUser.Error(), http.StatusInternalServerError)
		return

	}

	var response response.GeneralResponseNoData
	response.Message = "Delete user operation success"
	response.Status = http.StatusOK

	logging.InsertLog(r, constant.UserDeleteAccount, "", tokenParam, http.StatusOK, 4, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logout User func
func LogoutAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {

		logging.InsertLog(r, constant.UserDeleteAccount, errAunth.Error(), tokenParam, AunthStatus, 1, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	logoutUser, errLogout := module.LogoutAccountFromDB(GetUserIdAunth)

	if errLogout != nil {

		logging.InsertLog(r, constant.UserDeleteAccount, errLogout.Error(), tokenParam,  http.StatusInternalServerError, 1, 3)

		http.Error(w, errLogout.Error(), http.StatusInternalServerError)
		return
	}

	if !logoutUser {

		logging.InsertLog(r, constant.UserDeleteAccount, "Error when trying to logout", tokenParam,  http.StatusInternalServerError, 1, 3)

		http.Error(w, "Error when trying to logout", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Message = "Logout Succes"
	response.Status = http.StatusOK

	logging.InsertLog(r, constant.UserDeleteAccount, "", tokenParam,  http.StatusOK, 1, 4)


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Refresh JWT, place this on init
func RefrshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {

		logging.InsertLog(r, constant.RefreshToken, errAunth.Error(), tokenParam,  AunthStatus, 1, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	// Now check and refresh the token
	RefreshJWTandSession, errRefreshJWT := jwttoken.ReNewJWTandSession(tokenParam, GetUserIdAunth)

	if errRefreshJWT != nil {

		logging.InsertLog(r, constant.RefreshToken, errRefreshJWT.Error(), tokenParam,  http.StatusInternalServerError, 1, 3)

		http.Error(w, errRefreshJWT.Error(), http.StatusInternalServerError)
		return
	}

	var response response.ResponseUserLoginWithJWT
	response.Message = "Token Refreshed"
	response.Status = http.StatusOK
	response.JwtToken = RefreshJWTandSession

	logging.InsertLog(r, constant.RefreshToken, "", tokenParam,  http.StatusOK, 1, 4)

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
