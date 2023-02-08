package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"AzureWS/config"
	"crypto/md5"
	"encoding/hex"

	_ "github.com/lib/pq" // postgres golang driver
)

// Buku schema dari tabel Buku
// kita coba dengan jika datanya null
// jika return datanya ada yg null, silahkan pake NullString, contohnya dibawah
// Var       config.NullString `json:"var"`

type User struct {
	ID       *int64  `json:"id,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Token    *string `json:"token,omitempty"`
}


func AddUser(user User) int64 {

	db := config.CreateConnection()

	defer db.Close()

	// mengembalikan nilai id akan mengembalikan id dari user login yang dimasukkan ke db
	sqlStatement := `INSERT INTO user_login (username, password, token) VALUES ($1, $2, $3) RETURNING id`

	// id yang dimasukkan akan disimpan di id ini
	var id int64

	//membuat timestamp waktu sekarang tanpa format
	now := time.Now()

	//Generate Token menggunakan username, password, timestamp.now
	sum := md5.Sum([]byte(user.Password + user.Username + now.String()))
	tokenGenerated := hex.EncodeToString(sum[:])

	sumPswd := md5.Sum([]byte(user.Password))
	PasswordEncrpyted := hex.EncodeToString(sumPswd[:])

	// Scan function akan menyimpan insert id didalam id id
	err := db.QueryRow(sqlStatement, user.Username, PasswordEncrpyted, tokenGenerated).Scan(&id)

	if err != nil {
		log.Fatalf("Tidak Bisa mengeksekusi query. %v", err)
	}

	fmt.Printf("Insert data single record into user login %v", id)

	// return insert id
	return id
}


func GetAllUser() ([]User, error) {
	
	db := config.CreateConnection()

	defer db.Close()

	var users []User


	sqlStatement := `SELECT * FROM user_login`

	// mengeksekusi sql query
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("tidak bisa mengeksekusi query. %v", err)
	}

	// kita tutup eksekusi proses sql qeurynya
	defer rows.Close()

	// kita iterasi mengambil datanya
	for rows.Next() {
		var user User

		// kita ambil datanya dan unmarshal ke structnya
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Token)

		if err != nil {
			log.Fatalf("tidak bisa mengambil data semua user. %v", err)
		}

		// masukkan kedalam slice users
		users = append(users, user)

	}

	// return empty buku atau jika error
	return users, err
}


func GetSingleUser(id int64) (User, error) {
	// mengkoneksikan ke db postgres
	db := config.CreateConnection()

	// kita tutup koneksinya di akhir proses
	defer db.Close()

	var user User

	// buat sql query
	sqlStatement := `SELECT * FROM user_login WHERE id=$1`

	// eksekusi sql statement
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Token)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Tidak ada data yang dicari!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("tidak bisa mengambil data. %v", err)
	}

	return user, err
}

// update user in the DB
func UpdatePasswordUser(id int64, user User) int64 {

	// mengkoneksikan ke db postgres
	db := config.CreateConnection()

	// kita tutup koneksinya di akhir proses
	defer db.Close()

	// kita buat sql query create
	sqlStatement := `UPDATE user_login SET password=$2, WHERE id=$1 AND WHERE token=$3`

	// eksekusi sql statement
	res, err := db.Exec(sqlStatement, id, user.Password, user.Token)

	if err != nil {
		log.Fatalf("Tidak bisa mengeksekusi query ganti password. %v", err)
	}

	// cek berapa banyak row/data yang diupdate
	rowsAffected, err := res.RowsAffected()

	//kita cek
	if err != nil {
		log.Fatalf("Error ketika mengecheck rows/data yang diupdate. %v", err)
	}

	fmt.Printf("Total rows/record yang diupdate %v\n", rowsAffected)

	return rowsAffected
}

func RemoveUser(id int64, token string) int64 {

	// mengkoneksikan ke db postgres
	db := config.CreateConnection()

	// kita tutup koneksinya di akhir proses
	defer db.Close()

	// buat sql query
	sqlStatement := `DELETE FROM user_login WHERE id=$1 AND WHERE token=$2`

	// eksekusi sql statement
	res, err := db.Exec(sqlStatement, id, token)

	if err != nil {
		log.Fatalf("tidak bisa mengeksekusi query delete user. %v", err)
	}

	// cek berapa jumlah data/row yang di hapus
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("tidak bisa mencari data. %v", err)
	}

	fmt.Printf("Total data yang terhapus %v", rowsAffected)

	return rowsAffected
}
