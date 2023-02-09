package validation

import (
	"crypto/hmac"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"AzureWS/config"

)

func Validate(username, password string) (string, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT token FROM user_login WHERE username = $1 AND password = $2`

	// Execute the SQL statement.
	var token sql.NullString
	err := db.QueryRow(sqlStatement, username, password).Scan(&token)

	// If the user is not found, return an error.
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}

	// If there's an error in executing the SQL statement, return the error.
	if err != nil {
		log.Fatalf("Error executing the SQL statement: %v", err)
		return "", err
	}

	if token.Valid {
		return token.String, nil
	} else {
		return "", nil
	}
}

func ValidateTokenGetUuid(token string) (string, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT userId FROM user_login WHERE token = $1`

	// Execute the SQL statement.
	var uuid sql.NullString
	err := db.QueryRow(sqlStatement, token).Scan(&uuid)

	// If the user is not found, return an error.
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}

	// If there's an error in executing the SQL statement, return the error.
	if err != nil {
		log.Fatalf("Error executing the SQL statement: %v", err)
		return "", err
	}

	if uuid.Valid {
		return uuid.String, nil
	} else {
		return "", fmt.Errorf("INVALID USER ID")
	}
}

const salt = "AzureWsKey"

func ValidateUserPassword(enteredPassword string, storedPassword string) (bool, error) {
	// Hash the entered password using the same salt and function as used during registration
	h := sha512.New()
	h.Write([]byte(salt + enteredPassword))
	saltedPassword := h.Sum(nil)

	// Compare the salted password from the client with the stored salted password
	return hmac.Equal(saltedPassword, []byte(storedPassword)), nil
}

func ValidateGetStoredPassword (username string) (string, error){

	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT password FROM user_login WHERE username = $1`

	// Execute the SQL statement.
	var password sql.NullString

	err := db.QueryRow(sqlStatement, password).Scan(&password)
	
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%s","Password not found")
	}
	
	if err != nil {
		log.Fatalf("Error executing the SQL statement: %v", err)
		return "", err
	}

	return password.String, nil

}