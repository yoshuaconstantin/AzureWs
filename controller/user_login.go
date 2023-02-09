package controller

import (
	"encoding/json" // package untuk enkode dan mendekode json menjadi struct dan sebaliknya
	"fmt"
	"log"
	"net/http" // digunakan untuk mengakses objek permintaan dan respons dari api
	"strconv"  // package yang digunakan untuk mengubah string menjadi tipe int

	"github.com/gorilla/mux" // digunakan untuk mendapatkan parameter dari router
	_ "github.com/lib/pq"    // postgres golang driver

	"AzureWS/models" //models package dimana User didefinisikan
	"AzureWS/validation"
)

/*

 */

type responseUserLogin struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
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

// TambahUser
func InsrtNewUser(w http.ResponseWriter, r *http.Request) {

	// create an empty user of type models.User
	// kita buat empty User dengan tipe models.User
	var user models.User

	// decode data json request ke User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Tidak bisa mendecode dari request body.  %v", err)
	}

	// panggil modelsnya lalu insert User
	insertID := models.AddUser(user)

	// format response objectnya

	var response responseUserLogin
	response.Status = http.StatusOK
	response.Message = "Data user baru telah di tambahkan"
	response.ID = insertID

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

	username := loginData.Username
	password := loginData.Password

	//User should validate the password with stored password using salt and return boolean
	//change the whole system to validate first then continue
	storedPassword, err := validation.ValidateGetStoredPassword(username)

	PasswordValidation, errPassvalidate := validation.ValidateUserPassword(password, storedPassword)

	if errPassvalidate != nil {
		// kirim respon 400 kalau ada error

		var response ResponseAllError
		response.Status = http.StatusBadRequest
		response.Message = "Password not match"

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(response)
	}

	if PasswordValidation {

		token, err := validation.Validate(username, storedPassword)

		if err == nil {
			if token != "" {
				// Kirim respon token kalau tidak kosong
				var response ResponseUserToken
				response.Status = http.StatusOK
				response.Message = "Success"
				response.Token = token

				w.WriteHeader(http.StatusOK)

				json.NewEncoder(w).Encode(response)

			} else {
				// kirim respon 500 kalau token kosong

				var response ResponseUserToken
				response.Status = http.StatusInternalServerError
				response.Message = "No token"
				response.Token = token

				w.WriteHeader(http.StatusInternalServerError)

				json.NewEncoder(w).Encode(response)

			}
		} else {
			// kirim respon 400 kalau ada error

			var response ResponseAllError
			response.Status = http.StatusBadRequest
			response.Message = "User not found "

			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(response)
		}
	}
}

// Ambil semua data User
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

func UpdtUserPsswd(w http.ResponseWriter, r *http.Request) {

	// kita ambil request parameter idnya
	params := mux.Vars(r)

	// konversikan ke int yang sebelumnya adalah string
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Tidak bisa mengubah dari string ke int.  %v", err)
	}

	// buat variable User dengan type models.User
	var user models.User

	// decode json request ke variable User
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Tidak bisa decode request body.  %v", err)
	}

	// panggil updateUser untuk mengupdate data
	updatedRows := models.UpdatePasswordUser(int64(id), user)

	// ini adalah format message berupa string
	msg := fmt.Sprintf("User Password diupdate. Jumlah yang diupdate %v rows/record", updatedRows)

	// ini adalah format response message
	res := responseUserLogin{
		ID:      int64(id),
		Message: msg,
		Status:  http.StatusOK,
	}

	// kirim berupa response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// delete user
func DltUsr(w http.ResponseWriter, r *http.Request) {

	// kita ambil request parameter idnya
	params := mux.Vars(r)

	// konversikan ke int yang sebelumnya adalah string
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Tidak bisa mengubah dari string ke int.  %v", err)
	}

	token, ok := params["token"]

	if !ok {
		log.Fatalf("Data token tidak ada.  %v", ok)

		resFailed := responseUserLogin{
			ID:      int64(id),
			Message: "Data token tidak ada, silahkan masukan token",
			Status:  http.StatusBadRequest,
		}

		json.NewEncoder(w).Encode(resFailed)

	}

	// panggil fungsi hapusUser , dan convert int ke int64, masukin param token dari request
	deletedRows := models.RemoveUser(int64(id), token)

	// ini adalah format message berupa string
	msg := fmt.Sprintf("User Removed. Total data yang dihapus %v", deletedRows)

	// ini adalah format reponse message
	res := responseUserLogin{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}
