package module

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	jwttoken "AzureWS/JWTTOKEN"
	"AzureWS/config"
	Ct "AzureWS/globalvariable/constant"
	Gv "AzureWS/globalvariable/variable"
	"AzureWS/schemas/models"
	"AzureWS/session"
	"AzureWS/validation"

)

func CreateAccountToDB(user models.UserModel) (*models.TokenWithJwtModel, error) {
	var returnData models.TokenWithJwtModel

	db := config.CreateConnection()

	defer db.Close()

	//Validate if username is same
	usernameValidation, errUsername := validation.ValidateCreateNewUsername(user.Username)
	fmt.Printf("Username Validation Status: %v\n", usernameValidation)

	if errUsername != nil {
		fmt.Printf("\nCREATE USER - Masuk kedalam error\n")
		return nil, errUsername
	}

	//if Validation Username return true (means not having same username)
	if usernameValidation {

		// mengembalikan nilai id akan mengembalikan id dari user login yang dimasukkan ke db
		sqlStatement := `INSERT INTO user_login (user_id, username, password, token) VALUES ($1, $2, $3, $4) RETURNING id`

		// id yang dimasukkan akan disimpan di id ini
		var id int64

		//Generate UUID untuk userID menggunakan UUID generator
		userID, errUuid := uuid.NewRandom()
		if errUuid != nil {
			return nil, errUuid
		}

		var fixedUserId string = userID.String()

		//Encrypt password using salt
		//salt := []byte("AzureKey")
		hashedPassword, errhashed := validation.ValidatePasswordToEncrypt(user.Password)

		if errhashed != nil {
			return nil, errhashed
		}

		sum := md5.Sum([]byte(user.Username + user.Password + Gv.CurrentTime.String()))
		tokenGenerated := hex.EncodeToString(sum[:])

		createNewSession, errCreateNewSession := session.CreateNewSession(fixedUserId)

		if errCreateNewSession != nil {
			return nil, errCreateNewSession
		}

		if createNewSession {

			GenereateJWTToken, erroGenerateJwt := jwttoken.GenerateToken(fixedUserId)

			if erroGenerateJwt != nil {
				return nil, erroGenerateJwt
			}

			err := db.QueryRow(sqlStatement, fixedUserId, user.Username, hashedPassword, tokenGenerated).Scan(&id)

			if err != nil {
				return nil, err
			}

			fmt.Printf("\nCREATE USER - Insert data single record into user login %v\n", id)

			//Insert InitDashboards
			initDashboards, error := InitDashboardsDataSetToDB(fixedUserId)

			if error != nil {
				return nil, error
			}

			if !initDashboards {
				return nil, fmt.Errorf("%s", "\nCREATE USER - Failed to insert Init Dashboards Data\n")
			}

			//Insert InitProfile
			InitProfileData, errInitProfileData := InitUserProfileToDatabase(fixedUserId)

			if errInitProfileData != nil {
				return nil, errInitProfileData
			}

			if !InitProfileData {
				return nil, fmt.Errorf("%s", "\nINIT DATA PROFILE - Failed to insert Init Profile Data\n")
			}
			returnData.JWT = GenereateJWTToken
			returnData.Token = tokenGenerated

			return &returnData, nil

		} else {
			return nil, fmt.Errorf("%s", "\nCREATE SESSION - Failed to create new session\n")
		}

	} else {
		return nil, errUsername
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
func UpdatePasswordAccountFromDB(userId, password string) (bool, error) {

	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE user_login SET password=$1, WHERE user_id = $2`

	hashedPassword, errhashed := validation.ValidatePasswordToEncrypt(password)

	if errhashed != nil {
		return false, errhashed
	}

	res, err := db.Exec(sqlStatement, hashedPassword, userId)

	if err != nil {
		log.Fatalf("\nUpdate Password - Error : %v\n", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("\nUpdate Password - Error : %v\n", err)
	}

	if rowsAffected < 1 {
		return false, fmt.Errorf("%s","Cannot update password!")
	}

	return true, nil
}

/*
Note : To Delete All row where userId in single line use function CASCADE in SQl
Usage link *https://stackoverflow.com/questions/129265/cascade-delete-just-once*

And implement with your own style
*/
func RemoveAccountFromDB(userId string) (string, error) {

	db := config.CreateConnection()

	defer db.Close()

	sqlStatementUsrLgn := `DELETE FROM user_login WHERE user_id=$1`
	sqlStatementDshbrds := `DELETE FROM dashboards_data WHERE user_id=$1`
	sqlStatementUsrFdbck := `DELETE FROM user_feedback WHERE user_id=$1`
	sqlStatementUsrPrfl := `DELETE FROM user_profile WHERE user_id=$1`
	sqlStatementCmnyPst := `DELETE FROM community_post WHERE user_id = $1`
	sqlStatementCmnyCmnt := `DELETE FROM community_post_comment WHERE user_id = $1`
	sqlStatementCmnyLike := `DELETE FROM community_post_like WHERE user_id = $1`

	// eksekusi sql statement
	_, errDelUsrLgn := db.Exec(sqlStatementUsrLgn, userId)
	_, errDelDshbrds := db.Exec(sqlStatementDshbrds, userId)
	_, errDelUsrFdbck := db.Exec(sqlStatementUsrFdbck, userId)
	_, errDelUsrPrfl := db.Exec(sqlStatementUsrPrfl, userId)
	_, errDelCmnyPst := db.Exec(sqlStatementCmnyPst, userId)
	_, errDelCmnyCmnt := db.Exec(sqlStatementCmnyCmnt, userId)
	_, errDelCmnyLike := db.Exec(sqlStatementCmnyLike, userId)

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

	if errDelCmnyPst != nil {
		log.Fatalf("DELETE USER - OPERATION FAILED, REASON : %v", errDelCmnyPst)

		return "", errDelCmnyPst
	}

	if errDelCmnyCmnt != nil {
		log.Fatalf("DELETE USER - OPERATION FAILED, REASON : %v", errDelCmnyCmnt)

		return "", errDelCmnyCmnt
	}

	if errDelCmnyLike != nil {
		log.Fatalf("DELETE USER - OPERATION FAILED, REASON : %v", errDelCmnyLike)

		return "", errDelCmnyLike
	}

	return "DELETE USER - Operation succes", nil
}

// Logout user and delete the session
func LogoutAccountFromDB(userId string) (bool, error) {
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

func InsertTesting(password string) error {
	db := config.CreateConnection()

	defer db.Close()

	hashedPassword, errhashed := bcrypt.GenerateFromPassword(append(Ct.Salt, []byte(password)...), Ct.BcryptCost)

	if errhashed != nil {
		fmt.Printf("error generating password hash: %v", errhashed)
	}

	sqlStatement := `INSERT INTO testing (testencrypt) VALUES ($1)`

	_, err := db.Exec(sqlStatement, hashedPassword)

	if err != nil {
		return err
	}

	return nil
}

func CheckTesting(password string) (bool, error) {
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT testencrypt FROM testing WHERE id = 3`

	// Execute the SQL statement.
	var pswd sql.NullString

	err := db.QueryRow(sqlStatement).Scan(&pswd)

	if err == sql.ErrNoRows {
		return false, fmt.Errorf("%s", "Password not found")
	}

	if err != nil {
		log.Fatalf("Error executing the SQL statement: %v", err)
		return false, err
	}

	if err != nil {
		fmt.Printf("error generating password hash for validation: %v", err)
	}

	return bcrypt.CompareHashAndPassword([]byte(pswd.String), append(Ct.Salt, []byte(password)...)) == nil, nil
}
