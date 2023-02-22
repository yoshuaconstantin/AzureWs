package module

import (
	jwttoken "AzureWS/JWTTOKEN"
	"AzureWS/config"
	"AzureWS/schemas/models"
	"AzureWS/session"
	"AzureWS/validation"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // postgres golang driver
)

// jika return datanya ada yg null, silahkan pake NullString, contohnya dibawah
// Var       config.NullString `json:"var"`

func AddUserToDB(user models.UserModel) (string, error) {

	db := config.CreateConnection()

	defer db.Close()

	//Validate if username is same
	usernameValidation, errUsername := validation.ValidateCreateNewUsername(user.Username)
	fmt.Printf("Username Validation Status: %v\n", usernameValidation)

	if errUsername != nil {
		fmt.Printf("\nCREATE USER - Masuk kedalam error\n")
		return "", errUsername
	}

	//if Validation Username return true (means not having same username)
	if usernameValidation {

		// mengembalikan nilai id akan mengembalikan id dari user login yang dimasukkan ke db
		sqlStatement := `INSERT INTO user_login (user_id, username, password, token) VALUES ($1, $2, $3, $4) RETURNING id`

		// id yang dimasukkan akan disimpan di id ini
		var id int64

		//membuat timestamp waktu sekarang tanpa format
		now := time.Now()

		//Generate UUID untuk userID menggunakan UUID generator
		userID, errUuid := uuid.NewRandom()
		if errUuid != nil {
			fmt.Printf("\nCREATE USER - error generating UUID: %v\n", errUuid)
		}

		var fixedUserId string = userID.String()

		//Encrypt password using salt
		//salt := []byte("AzureKey")
		hashedPassword, errhashed := validation.ValidatePasswordToEncrypt(user.Password)

		fmt.Printf("\nCREATE USER - Generated Password Salt %v\n", hashedPassword)

		if errhashed != nil {
			fmt.Printf("\nCREATE USER - error generating password hash: %v\n", errhashed)
		}

		//Generate Token menggunakan username, password, timestamp.now
		sum := md5.Sum([]byte(user.Password + user.Username + now.String()))
		tokenGenerated := hex.EncodeToString(sum[:])

		createNewSession, errCreateNewSession := session.CreateNewSession(fixedUserId)

		if errCreateNewSession != nil {
			return "", errCreateNewSession
		}

		if createNewSession {

			GenereateJWTToken, erroGenerateJwt := jwttoken.GenerateToken(fixedUserId)

			if erroGenerateJwt != nil {
				return "", erroGenerateJwt
			}

			err := db.QueryRow(sqlStatement, fixedUserId, user.Username, hashedPassword, tokenGenerated).Scan(&id)

			if err != nil {
				log.Fatalf("\nCREATE USER - Tidak Bisa mengeksekusi query. %v\n", err)
			}

			fmt.Printf("\nCREATE USER - Insert data single record into user login %v\n", id)

			//Insert InitDashboards
			initDashboards, error := InitDashboardsDataSetToDB(fixedUserId)

			if error != nil {
				return "", error
			}

			if !initDashboards {
				return "", fmt.Errorf("%s", "\nCREATE USER - Failed to insert Init Dashboards Data\n")
			}

			//Insert InitProfile
			InitProfileData, errInitProfileData := InitUserProfileToDatabase(fixedUserId)

			if errInitProfileData != nil {
				return "", errInitProfileData
			}

			if !InitProfileData {
				return "", fmt.Errorf("%s", "\nINIT DATA PROFILE - Failed to insert Init Profile Data\n")
			}

			return GenereateJWTToken, nil

		} else {
			return "", fmt.Errorf("%s", "\nCREATE SESSION - Failed to create new session\n")
		}

	} else {
		return "", errUsername
	}
}

func GetAllUser() ([]models.UserModel, error) {

	db := config.CreateConnection()

	defer db.Close()

	var users []models.UserModel

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
		var user models.UserModel

		// kita ambil datanya dan unmarshal ke structnya
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.UserId)

		if err != nil {
			log.Fatalf("tidak bisa mengambil data semua user. %v", err)
		}

		// masukkan kedalam slice users
		users = append(users, user)

	}

	// return empty buku atau jika error
	return users, err
}

func GetSingleUser(id int64) (models.UserModel, error) {
	// mengkoneksikan ke db postgres
	db := config.CreateConnection()

	// kita tutup koneksinya di akhir proses
	defer db.Close()

	var user models.UserModel

	// buat sql query
	sqlStatement := `SELECT * FROM user_login WHERE id=$1`

	// eksekusi sql statement
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.UserId)

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
func UpdatePasswordUserFromDB(userId, password string) (int64, error) {

	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE user_login SET password=$1, WHERE user_id = $2`

	res, err := db.Exec(sqlStatement, password, userId)

	if err != nil {
		log.Fatalf("\nUpdate Password - Tidak bisa mengeksekusi query ganti password. %v\n", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("\nUpdate Password - Error ketika mengecheck rows/data yang diupdate. %v\n", err)
	}

	fmt.Printf("\nUpdate Password - State succes\n")

	return rowsAffected, nil
}

/*
Note : To Delete All row where userId in single line use function CASCADE in SQl
Usage link *https://stackoverflow.com/questions/129265/cascade-delete-just-once*

And implement with your own style
*/
func RemoveUserDB(userId string) (string, error) {

	db := config.CreateConnection()

	defer db.Close()

	sqlStatementUsrLgn := `DELETE FROM user_login WHERE user_id=$1`
	sqlStatementDshbrds := `DELETE FROM dashboards_data WHERE user_id=$1`
	sqlStatementUsrFdbck := `DELETE FROM user_feedback WHERE user_id=$1`
	sqlStatementUsrPrfl := `DELETE FROM user_profile WHERE user_id=$1`

	// eksekusi sql statement
	_, errDelUsrLgn := db.Exec(sqlStatementUsrLgn, userId)
	_, errDelDshbrds := db.Exec(sqlStatementDshbrds, userId)
	_, errDelUsrFdbck := db.Exec(sqlStatementUsrFdbck, userId)
	_, errDelUsrPrfl := db.Exec(sqlStatementUsrPrfl, userId)

	if errDelUsrLgn != nil {
		log.Fatalf("DELETE USER - OPERATION FAILED, REASON : %v", errDelUsrLgn)

		return "", errDelUsrLgn
	}

	if errDelDshbrds != nil {
		log.Fatalf("DELETE USER - OPERATION FAILED, REASON : %v", errDelDshbrds)

		return "", errDelDshbrds
	}

	if errDelUsrFdbck != nil {
		log.Fatalf("DELETE USER - OPERATION FAILED, REASON : %v", errDelUsrFdbck)

		return "", errDelUsrFdbck
	}

	if errDelUsrPrfl != nil {
		log.Fatalf("DELETE USER - OPERATION FAILED, REASON : %v", errDelUsrPrfl)

		return "", errDelUsrPrfl
	}

	return "DELETE USER - Operation succes", nil
}

// Logout user and delete the session
func LogoutUser(userId string) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE user_session SET session_id = '', is_active = 'false' WHERE user_id = $1`

	_, err := db.Exec(sqlStatement, userId)

	if err != nil {
		log.Fatalf("\nLOGOUT USER - Cannot execute command : %v\n", err)
		return false, err
	}

	return true, nil
}