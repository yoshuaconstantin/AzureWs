package validation

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"AzureWS/config"
	Gv "AzureWS/globalvariable/variable"
	Ct "AzureWS/globalvariable/constant"

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

func ValidateGenerateNewToken(username, password string) (string, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `UPDATE user_login SET token = $1 WHERE username = $2`

	
	sum := md5.Sum([]byte(password + username + Gv.CurrentTime.String()))
	tokenGenerated := hex.EncodeToString(sum[:])

	_, errExec := db.Exec(sqlStatement, tokenGenerated, username)

	if errExec != nil {
		return "", errExec
	}

	return tokenGenerated, nil
}

// Validate users token to get user id
func ValidateTokenGetUuid(token string) (string, error) {
	// Connect to the database.
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT user_id FROM user_login WHERE token = $1`

	//token print
	fmt.Printf("\nToken Validation - Token Entered %v\n", token)

	// Execute the SQL statement.
	var uuid sql.NullString
	err := db.QueryRow(sqlStatement, token).Scan(&uuid)

	// If the user is not found, return an error.
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%s", "Token Validation - user not found")
	}

	// If there's an error in executing the SQL statement, return the error.
	if err != nil {
		log.Fatalf("Error executing the SQL statement: %v", err)
		return "", err
	}

	if uuid.Valid {
		return uuid.String, nil
	} else {
		return "", fmt.Errorf("%s", "Token Validation - Invalid Token")
	}
}

// Validate users username to get stored password
func ValidateGetStoredPassword(username string) (string, error) {

	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT password FROM user_login WHERE username = $1`

	// Execute the SQL statement.
	var password sql.NullString

	err := db.QueryRow(sqlStatement, username).Scan(&password)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%s", "Password not found")
	}

	if err != nil {
		log.Fatalf("Error executing the SQL statement: %v", err)
		return "", err
	}

	return password.String, nil
}

// Validate username when user create an account
func ValidateCreateNewUsername(username string) (bool, error) {
	db := config.CreateConnection()

	// Close the connection at the end of the process.
	defer db.Close()

	// Create a SQL query to retrieve the token based on the username and password.
	sqlStatement := `SELECT username FROM user_login WHERE username = $1`

	// Execute the SQL statement.

	var result string
	err := db.QueryRow(sqlStatement, username).Scan(&result)

	if err == sql.ErrNoRows {
		fmt.Printf("\nVALIDATE USERNAME - No rows\n")
		return true, nil
	}

	if err != nil {
		log.Fatalf("\nVALIDATE USERNAME - Error executing the SQL statement: %v\n", err)
		return false, err
	}

	if result == username {
		return false, nil
	} else {
		return true, nil
	}
}

// Validate users password before login into account
func ValidateUserPassword(enteredPassword string, storedPassword string) (bool, error) {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), append(Ct.Salt, []byte(enteredPassword)...)) == nil, nil
}

// Validate user password to get encrypted password
func ValidatePasswordToEncrypt(password string) (string, error) {
	hashedPassword, errhashed := bcrypt.GenerateFromPassword(append(Ct.Salt, []byte(password)...), Ct.BcryptCost)

	if errhashed != nil {
		fmt.Printf("error generating password hash: %v", errhashed)
	}
	return string(hashedPassword), nil
}
