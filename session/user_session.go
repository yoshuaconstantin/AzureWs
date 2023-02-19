package session

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"
	"AzureWS/config"

	_ "github.com/lib/pq"
	//"golang.org/x/crypto/bcrypt"
)

// Check session when login, should always return true to add expired
func CheckSessionLogin(userId string) (bool, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT expired FROM user_session WHERE user_id = $1`

	// Execute the SQL statement.
	var isExpired sql.NullString
	err := db.QueryRow(sqlStatement, userId).Scan(&isExpired)

	// If the user is not found, return an error.
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("%s", "\nSESSION CHECKING - Session not found\n")
	}

	// If there's an error in executing the SQL statement, return the error.
	if err != nil {
		log.Fatalf("\nSESSION CHECKING - Error executing the SQL statement: %v\n", err)
		return false, err
	}

	currentTime := time.Now()
	expiry := currentTime.Add(time.Hour * 24 * 3)
	expiryStr := expiry.Format("2006-01-02 15:04")

	if isExpired.Valid {

		// Add Duration of expired when login and session not expired
		AddExpired := `UPDATE user_login SET expired = $1, WHERE user_id = $2 AND is_active = 'true'`

		_, err := db.Exec(AddExpired, expiryStr, userId)

		if err != nil {
			return false, fmt.Errorf("%s", "\nSESSION CHECKING - Cannot update new expired date\n")
		}

		return true, nil
	} else {
		return false, fmt.Errorf("%s", "\nSESSION CHECKING - Session is empty\n")
	}
}

// Check Session after login
func CheckSessionInside(userId string) (bool, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT expired FROM user_session WHERE user_id = $1 AND is_active = 'true'`

	// Execute the SQL statement.
	var isExpired sql.NullString
	err := db.QueryRow(sqlStatement, userId).Scan(&isExpired)

	// If the user is not found, return an error.
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("%s", "\nSESSION CHECKING - Session not found\n")
	}

	// If there's an error in executing the SQL statement, return the error.
	if err != nil {
		log.Fatalf("\nSESSION CHECKING - Error executing the SQL statement: %v\n", err)
		return false, err
	}

	currentTime := time.Now()
	expiryStr := currentTime.Format("2006-01-02 15:04")

	if isExpired.Valid {
		if expiryStr > isExpired.String {
			return false, fmt.Errorf("%s", "\nSESSION CHECKING - Session Expired, log in again.\n")
		} else {
			return true, nil
		}
	} else {
		return false, nil
	}
}

// Create new session when users create new account
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

// Generate sessionId random string
func GenerateSessionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
