package models

import (
	"fmt"
	"log"

	_ "github.com/lib/pq" // postgres golang driver

	"AzureWS/config"
	"AzureWS/session"
	"AzureWS/validation"

)

type UserProfileData struct {
	UserId       *int    `json:"user_id,omitempty"`
	Nickname     *string `json:"nickname,omitempty"`
	Age          *int    `json:"age,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	UserImage    *string `json:"user_image,omitempty"`
	CreatedSince *string `json:"created_since,omitempty"`
	Token        *string `json:"token,omitempty"`
}


// Can be used to insert and update profile data
func InsertUserProfileToDatabase(userData UserProfileData) (string, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO user_profile (user_id, nickname, age, gender, user_image, created_since) VALUES ($1, $2, $3, $4, $5, $6)`

	SessionValidation, errSessionCheck := session.CheckSession(*userData.Token)

	if errSessionCheck != nil {
		return "", errSessionCheck
	}

	if SessionValidation {

		GetUserID, errGetUuid := validation.ValidateTokenGetUuid(*userData.Token)

		if errGetUuid != nil {
			return "", errGetUuid
		}

		_, err := db.Exec(sqlStatement, GetUserID, userData.Nickname, userData.Age, userData.Gender, userData.UserImage, userData.CreatedSince)

		if err != nil {
			log.Fatalf("\nINSERT USER PROFILE - Cannot execute command : %v\n", err)
		}

		return "Profile Updated", nil

	} else {
		return "", fmt.Errorf("%s", "SESSION VALIDATION - Session Expired")
	}
}
