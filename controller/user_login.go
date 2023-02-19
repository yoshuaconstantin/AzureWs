package controller

import (
	"encoding/json" // package untuk enkode dan mendekode json menjadi struct dan sebaliknya
	"fmt"
	"log"
	"net/http" // digunakan untuk mengakses objek permintaan dan respons dari api
	"strconv"  // package yang digunakan untuk mengubah string menjadi tipe int

	"github.com/gorilla/mux" // digunakan untuk mendapatkan parameter dari router
	_ "github.com/lib/pq"    // postgres golang driver

	jwttoken "AzureWS/JWTTOKEN"
	"AzureWS/models" //models package dimana User didefinisikan
	"AzureWS/session"
	"AzureWS/validation"

)

/*

 */

type responseUserLogin struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

type responseUserLoginWithJWT struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
	JwtToken string `json:"jtwToken,omitempty"`
}

type ResponseUserLogin struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []models.User `json:"data"`
}

type ResponseUserToken struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ResponseAllError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordData struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type TokenData struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func InsrtNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// create an empty user of type models.User
	// kita buat empty User dengan tipe models.User
	var user models.User

	// decode data json request ke User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("\nCREATE NEW USER - Cannot get request body : %v\n", err)
	}

	// panggil modelsnya lalu insert User
	InsertAndGetJwt, errInsert := models.AddUser(user)

	if errInsert != nil {
		var response ResponseAllError
		response.Status = http.StatusConflict
		response.Message = "Username already registered"

		// kirim response
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	} else {
		if InsertAndGetJwt == "" {
			var response ResponseAllError
			response.Status = http.StatusConflict
			response.Message = "Username already registered"

			// kirim response
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		} else {
			var response responseUserLoginWithJWT
			response.Status = http.StatusOK
			response.Message = "Data user baru telah di tambahkan"
			response.JwtToken = InsertAndGetJwt

			// kirim response
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	}
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
	user, err := models.GetSingleUser(int64(id))

	if err != nil {
		log.Fatalf("Tidak bisa mengambil data User. %v", err)
	}

	// kirim response
	json.NewEncoder(w).Encode(user)
}

// Login user
func LoginUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var loginData LoginData
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		var response ResponseAllError
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// If the authorization header is empty, return an error
		var response ResponseAllError
		response.Status = http.StatusBadRequest
		response.Message = "Missing authorization header"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	username := loginData.Username
	password := loginData.Password

	//User should validate the password with stored password using salt and return boolean
	//change the whole system to validate first then continue
	storedPassword, errStrdPswd := validation.ValidateGetStoredPassword(username)

	if errStrdPswd != nil {
		var response ResponseAllError
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("\nLOGIN USER - Error: %s\n", errStrdPswd.Error())

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	PasswordValidation, errPassvalidate := validation.ValidateUserPassword(password, storedPassword)

	fmt.Printf("Password Validation = %v\n", PasswordValidation)
	if errPassvalidate != nil {
		// kirim respon 400 kalau ada error

		var response ResponseAllError
		response.Status = http.StatusBadRequest
		response.Message = "Password not match"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if PasswordValidation {

		token, err := validation.Validate(username, storedPassword)

		if err == nil {
			if token != "" {

				GetUserID, errGetUuid := validation.ValidateTokenGetUuid(token)

				if errGetUuid != nil {
					var response responseUserProfile
					response.Status = http.StatusUnauthorized
					response.Message = errGetUuid.Error()
			
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(response)
					return
				}

				CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(GetUserID)

				if erroCheckJWt != nil {
					var response responseUserProfile
					response.Status = http.StatusUnauthorized
					response.Message = erroCheckJWt.Error()
			
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(response)
					return
				}

				if !CheckJwtTokenValidation {
					var response responseUserProfile
					response.Status = http.StatusUnauthorized
					response.Message = "Unauthorized user"
			
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(response)
					return
				} 

				checkLoginSession, errAddSession := session.CheckSessionLogin(GetUserID)

				if errAddSession != nil {

					var response ResponseUserToken
					response.Status = http.StatusInternalServerError
					response.Message = errAddSession.Error()
					response.Token = token
	
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				if !checkLoginSession {
					var response responseUserProfile
					response.Status = http.StatusUnauthorized
					response.Message = "Contact Dev to fix this"
			
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(response)
					return
				}

				fmt.Printf("\nPASSWORD VALIDATION - User Token = %v\n", token)
				// Kirim respon token kalau tidak kosong
				var response ResponseUserToken
				response.Status = http.StatusOK
				response.Message = "Success"
				response.Token = token

				w.WriteHeader(http.StatusOK)

				json.NewEncoder(w).Encode(response)
				return

			} else {
				// kirim respon 500 kalau token kosong

				var response ResponseUserToken
				response.Status = http.StatusInternalServerError
				response.Message = "No token"
				response.Token = token

				w.WriteHeader(http.StatusInternalServerError)

				json.NewEncoder(w).Encode(response)
				return

			}
		} else {

			var response ResponseAllError
			response.Status = http.StatusBadRequest
			response.Message = "User not found "

			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(response)
			return
		}
	} else {

		var response ResponseAllError
		response.Status = http.StatusBadRequest
		response.Message = "Password validating failed\n"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func GetAllUsr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Call the GetAllUser method from the models package
	users, err := models.GetAllUser()

	if err != nil {
		log.Fatalf("Unable to retrieve data. %v", err)

		var response ResponseAllError
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving data"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	var response ResponseUserLogin
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = users

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update User Password func
func UpdtUserPsswd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var updatePswdModel PasswordData

	err := json.NewDecoder(r.Body).Decode(&updatePswdModel)

	if err != nil {
		res := responseUserLogin{
			Message: "Cannot get request body",
			Status:  http.StatusBadRequest,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	userId, errGetUuid := validation.ValidateTokenGetUuid(updatePswdModel.Token)

	if errGetUuid != nil {
		log.Fatalf("Unable to retrieve UserId. %v", errGetUuid)

		var response responseDashboards
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving UserId"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	SessionValidation, errSessionCheck := session.CheckSessionInside(userId)

	if errSessionCheck != nil {
		var response responseUserProfile
		response.Status = http.StatusForbidden
		response.Message = errSessionCheck.Error()

		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !SessionValidation {
		var response responseUserProfile
		response.Status = http.StatusUnauthorized
		response.Message = "Session Expired"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, errUpdatePswd := models.UpdatePasswordUser(updatePswdModel.Token, updatePswdModel.Password)

	if errUpdatePswd != nil {
		res := responseUserLogin{
			Message: errUpdatePswd.Error(),
			Status:  http.StatusOK,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	res := responseUserLogin{
		Message: "Update password succes",
		Status:  http.StatusOK,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// Delete User func
func DltUsr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var tokenModel TokenData

	err := json.NewDecoder(r.Body).Decode(&tokenModel)

	if err != nil {
		res := responseUserLogin{
			Message: "Cannot get request body",
			Status:  http.StatusBadRequest,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	userId, errGetUuid := validation.ValidateTokenGetUuid(tokenModel.Token)
	
	if errGetUuid != nil {
		log.Fatalf("Unable to retrieve UserId. %v", errGetUuid)

		var response responseDashboards
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving UserId"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	SessionValidation, errSessionCheck := session.CheckSessionInside(userId)

	if errSessionCheck != nil {
		var response responseUserProfile
		response.Status = http.StatusForbidden
		response.Message = errSessionCheck.Error()

		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !SessionValidation {
		var response responseUserProfile
		response.Status = http.StatusUnauthorized
		response.Message = "Session Expired"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, errDeleteUser := models.RemoveUser(tokenModel.Token)

	if errDeleteUser != nil {

		var response responseUserLogin
		response.Message = "Delete user operation failed"
		response.Status = http.StatusBadRequest

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return

	}

	var response responseUserLogin
	response.Message = "Delete user operation success"
	response.Status = http.StatusOK

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logout User func
func LgoutUsr (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var tokenModel TokenData

	err := json.NewDecoder(r.Body).Decode(&tokenModel)

	if err != nil {
		res := responseUserLogin{
			Message: "Cannot get request body",
			Status:  http.StatusBadRequest,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	userId, errGetUuid := validation.ValidateTokenGetUuid(tokenModel.Token)
	
	if errGetUuid != nil {
		log.Fatalf("Unable to retrieve UserId. %v", errGetUuid)

		var response responseUserProfile
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving UserId"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	SessionValidation, errSessionCheck := session.CheckSessionInside(userId)

	if errSessionCheck != nil {
		var response responseUserProfile
		response.Status = http.StatusForbidden
		response.Message = errSessionCheck.Error()

		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !SessionValidation {
		var response responseUserProfile
		response.Status = http.StatusUnauthorized
		response.Message = "Session Expired"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	logoutUser, errLogout := models.LogoutUser(userId)

	if errLogout != nil {
		var response responseUserProfile
		response.Status = http.StatusInternalServerError
		response.Message = errLogout.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !logoutUser {
		var response responseUserProfile
		response.Status = http.StatusInternalServerError
		response.Message = "Error when trying to logout"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response responseUserLogin
	response.Message = "Logout Succes"
	response.Status = http.StatusOK

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}