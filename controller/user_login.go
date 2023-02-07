package controller

import (
	"encoding/json" // package untuk enkode dan mendekode json menjadi struct dan sebaliknya
	"fmt"
	"strconv" // package yang digunakan untuk mengubah string menjadi tipe int

	"log"
	"net/http" // digunakan untuk mengakses objek permintaan dan respons dari api

	"go-postgres-crud/models" //models package dimana User didefinisikan

	"github.com/gorilla/mux" // digunakan untuk mendapatkan parameter dari router
	_ "github.com/lib/pq"    // postgres golang driver
)

type responseUserLogin struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResponseUserLogin struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []models.User `json:"data"`
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
	res := responseUserLogin{
		ID:      insertID,
		Message: "Data user baru telah di tambahkan",
	}

	// kirim response
	json.NewEncoder(w).Encode(res)
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

// Ambil semua data User
func GetAllUsr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// memanggil models GetAllUser
	users, err := models.GetAllUser()

	if err != nil {
		log.Fatalf("Tidak bisa mengambil data. %v", err)
	}

	var response ResponseUserLogin
	response.Status = 1
	response.Message = "Success"
	response.Data = users

	// kirim semua response
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
	}

	// kirim berupa response
	json.NewEncoder(w).Encode(res)
}

func DltUsr(w http.ResponseWriter, r *http.Request) {

	// kita ambil request parameter idnya
	params := mux.Vars(r)

	// konversikan ke int yang sebelumnya adalah string
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Tidak bisa mengubah dari string ke int.  %v", err)
	}

	// panggil fungsi hapusUser , dan convert int ke int64, masukin param token dari request
	deletedRows := models.RemoveUser(int64(id), params["token"])

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
