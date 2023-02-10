package validation

import (
	"database/sql"
	"fmt"
	"log"
	"AzureWS/config"

	_ "github.com/lib/pq"
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
	sqlStatement := `SELECT user_id FROM user_login WHERE token = $1`

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


func ValidateGetStoredPassword (username string) (string, error){

	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT password FROM user_login WHERE username = $1`

	// Execute the SQL statement.
	var password sql.NullString

	err := db.QueryRow(sqlStatement, username).Scan(&password)
	
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%s","Password not found")
	}
	
	if err != nil {
		log.Fatalf("Error executing the SQL statement: %v", err)
		return "", err
	}

	return password.String, nil

}

func ValidateCreateNewUsername (username string) (bool, error){
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT username FROM user_login WHERE username = $1`

	// Execute the SQL statement.
	
	var result string
	err := db.QueryRow(sqlStatement, username).Scan(&result)
	
	if err == sql.ErrNoRows {
		fmt.Printf("Masuk kedalam tidak ada row = sukses\n")	
		return true, nil
	}
	
	if err != nil {
		fmt.Printf("Masuk kedalam error false")	
		log.Fatalf("Error executing the SQL statement: %v", err)
		return false, err
	}

	return true, nil
}