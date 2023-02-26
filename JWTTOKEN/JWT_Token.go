package jwttoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"AzureWS/session"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

/*
NOTE : Secret key should be stored in safe location. how to store and get the secret key something like this

create an example file in your local / vps setup then set the environment variable using terminal
*export SECRET_FILE_PATH=/path/to/your/secret/file*

then call using this code to get the secret key string
*keyBytes, err := ioutil.ReadFile(filePath)*

the func call should be like this
*secretKey, err := readSecretKey(filePath)*

*/

// Generate JWT with secret key and userId
func GenerateToken(userId string) (string, error) {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(), // Set expiration time to 3 days from now
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("testingkey"))
}

// Verify JWT using token string from header authorization bearer
func VerifyToken(tokenString string) (bool, error) {
	// Generate a random secret key
	secretKey := []byte("testingkey")

	// Parse the token with the secret key
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("%s %v", "Unexpected signing method: ", token.Header["alg"])
		}

		// Return the secret key
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}

	// Extract the user ID from the token
	_, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return false, fmt.Errorf("%s", "Unexpected token claims")
	}

	return true, nil
}

// Refresh JWT and return the latest token
func ReNewJWTandSession(tokenString, userId string) (string, error) {
	verifyToken, errVerify := VerifyToken(tokenString)

	if errVerify != nil {
		return "", errVerify
	}

	if !verifyToken {
		return "", fmt.Errorf("%s", "Unauthorized User")
	}

	generateNewToken, errGenerateToken := GenerateToken(userId)

	if errGenerateToken != nil {
		return "", errGenerateToken
	}

	_, errS := session.ReNewSessionLogin(userId)

	if errS != nil {
		return "", errS
	}

	return generateNewToken, nil
}
