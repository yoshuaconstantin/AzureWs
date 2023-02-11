package session

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	//"golang.org/x/crypto/bcrypt"

	"AzureWS/config"

)

func CheckSession(userId string) (bool, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT is_active FROM user_session WHERE user_id = $1`

	// Execute the SQL statement.
	var isActive sql.NullString
	err := db.QueryRow(sqlStatement, userId).Scan(&isActive)

	// If the user is not found, return an error.
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("%s", "\nSESSION CHECKING - Session not found\n")
	}

	// If there's an error in executing the SQL statement, return the error.
	if err != nil {
		log.Fatalf("\nSESSION CHECKING - Error executing the SQL statement: %v\n", err)
		return false, err
	}

	if isActive.Valid {
		if isActive.String == "true" {
			return true, nil
		} else {
			return false, fmt.Errorf("%s", "\nSESSION CHECKING - User do not have session active, re login!\n")
		}
	} else {
		return false, nil
	}
}

func CreateNewSession(userId string) (bool, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	sqlStatement := `INSERT INTO user_session (user_id, session_id, expired, is_active) VALUES ($1, $2, $3, 'true')`

	sessionID, errCreateSession := GenerateSessionID()
	if errCreateSession != nil {
		fmt.Println(" \nSESSION CREATE - Error generating session ID:\n", errCreateSession)

	}

	//Generate expired session today + next 3 days.
	currentTime := time.Now()
	expiry := currentTime.Add(time.Hour * 24 * 3)
	expiryStr := expiry.Format("2006-01-02 15:04")

	// Execute the SQL statement.
	_, err := db.Exec(sqlStatement, userId, sessionID, expiryStr)

	// If there's an error in executing the SQL statement, return the error.
	if err != nil {
		log.Fatalf("\nSESSION CREATE - Error executing the SQL statement: %v\n", err)
		return false, err
	}

	return true, nil
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
