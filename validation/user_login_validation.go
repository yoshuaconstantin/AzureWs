package validation

import (
	"log"
	"database/sql"
	"AzureWS/config"
	"fmt"

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