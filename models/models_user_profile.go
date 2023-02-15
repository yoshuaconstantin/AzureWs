package models

import (
	"log"

	_ "github.com/lib/pq" // postgres golang driver

	"AzureWS/config"

)

type UserProfileData struct {
	UserId       *int    `json:"user_id,omitempty"`
	Nickname     *string `json:"nickname,omitempty"`
	Age          *int    `json:"age,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	ImageUrl     *string `json:"image_url,omitempty"`
	CreatedSince *string `json:"created_since,omitempty"`
	Token        *string `json:"token,omitempty"`
}

type GetUserProfileData struct {
	Nickname     *string `json:"nickname,omitempty"`
	Age          *int    `json:"age,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	ImageUrl     *string `json:"image_url,omitempty"`
	CreatedSince *string `json:"created_since,omitempty"`
}

// Can be used to insert and update profile data
func InsertUserProfileToDatabase(userData UserProfileData, userId string) (string, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO user_profile (user_id, nickname, age, gender, image_url, created_since) VALUES ($1, $2, $3, $4, $5, $6)`

		_, err := db.Exec(sqlStatement, userId, userData.Nickname, userData.Age, userData.Gender, userData.ImageUrl, userData.CreatedSince)

		if err != nil {
			log.Fatalf("\nINSERT USER PROFILE - Cannot execute command : %v\n", err)
		}

		return "Profile Updated", nil

	
}

/*
Get User profile data, for image it should just return the file name only
then hit the img API to get the img url, store it to user model.

Task : Create endpoint for Hit image with GET method and Json filePath

Step : Hit GetUserProfileData -> return:model -> hit GetImageData (request: userId, filename)
return:imgUrl -> user should use Image.Network (Flutter)
*/
func GetUserProfileDataFromDatabase(userId string) ([]GetUserProfileData, error) {
	db := config.CreateConnection()

	defer db.Close()

	var profileData []GetUserProfileData

	sqlStatement := `SELECT nickname,age,gender,image_url,created_since FROM user_profile where user_id = $1`

	rows, err := db.Query(sqlStatement, userId)

	if err != nil {
		log.Fatalf("GET USER DATA PROFILE - Cannot exec query, error : %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var getUserProfileData GetUserProfileData

		err = rows.Scan(&getUserProfileData.Nickname, &getUserProfileData.Age, &getUserProfileData.Gender, &getUserProfileData.ImageUrl, &getUserProfileData.CreatedSince)

		if err != nil {
			log.Fatalf("GET USER DATA PROFILE - Error : %v", err)
		}

		profileData = append(profileData, getUserProfileData)
	}

	return profileData, nil
}
