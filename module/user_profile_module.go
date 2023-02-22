package module

import (
	"AzureWS/config"
	"AzureWS/schemas/models"
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Init insert user profile with user id and the rest string null
func InitUserProfileToDatabase(userId string) (bool, error) {

	db := config.CreateConnection()

	defer db.Close()

	// Created Date for user
	currentTime := time.Now()
	CreatedDate := currentTime.Format("2006-01-02 15:04")

	sqlStatement := `INSERT INTO user_profile (user_id, nickname, age, gender, nation, image_url, created_since) VALUES ($1, '', '','','','',$2)`

	_, err := db.Exec(sqlStatement, userId, CreatedDate)

	if err != nil {
		log.Fatalf("\nINSERT USER PROFILE - Cannot execute command : %v\n", err)
		return false, err
	}

	return true, nil
}

// Can be used to insert and update profile data
func UpdateUserProfileToDatabase(userData models.UserProfileDataModel, userId string) (string, error) {
	db := config.CreateConnection()

	defer db.Close()

	fmt.Println(userData)
	fmt.Println(userId)

	sqlStatement := `UPDATE user_profile SET nickname = $1, age = $2, gender = $3, image_url = $4, nation = $5 WHERE user_id = $6`

	_, err := db.Exec(sqlStatement, userData.Nickname, userData.Age, userData.Gender, userData.ImageUrl, userData.Nation,userId)

	if err != nil {
		log.Fatalf("\nUPDATE USER PROFILE - Cannot execute command : %v\n", err)
		return "", err
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
func GetUserProfileDataFromDatabase(userId string) ([]models.GetUserProfileDataModel, error) {
	db := config.CreateConnection()

	defer db.Close()

	var profileData []models.GetUserProfileDataModel

	sqlStatement := `SELECT nickname,age,gender,image_url,nation,created_since FROM user_profile where user_id = $1`

	rows, err := db.Query(sqlStatement, userId)

	if err != nil {
		log.Fatalf("GET USER DATA PROFILE - Cannot exec query, error : %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var getUserProfileData models.GetUserProfileDataModel

		err = rows.Scan(&getUserProfileData.Nickname, &getUserProfileData.Age, &getUserProfileData.Gender, &getUserProfileData.ImageUrl, &getUserProfileData.CreatedSince)

		if err != nil {
			log.Fatalf("GET USER DATA PROFILE - Error : %v", err)
		}

		profileData = append(profileData, getUserProfileData)
	}

	return profileData, nil
}

// Upload photo and update the database Image Url
func UploadUserProfilePhotoBool(userId string, byteImage []byte) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE user_profile SET image_url = $1 WHERE user_id = $2`

	ImageUrl, errGetImgUrl := ConvertByteToImgString(byteImage, userId)

	if errGetImgUrl != nil {
		return false, errGetImgUrl
	}

	_, err := db.Exec(sqlStatement, ImageUrl, userId)

	if err != nil {
		log.Fatalf("\nUPLOAD USER PROFILE IMAGE - Cannot execute command : %v\n", err)
		return false, err
	}

	return true, nil
}

// Process Byte To Img and save it with return of ImageUrl
func ConvertByteToImgString(byteImage []byte, userId string) (string, error) {

	img, format, err := image.Decode(bytes.NewReader(byteImage))
	if err != nil {
		log.Println(err)

		return "", err
	}

	path := filepath.Join("FileAzure", "Data", "Image", userId, "profile")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)

		return "", err
	}

	// Choose a unique filename for the image
	filename := fmt.Sprintf("%d.%s", time.Now().Unix(), format)
	filepath := filepath.Join(path, filename)

	// Create the file and write the image to it
	f, err := os.Create(filepath)
	if err != nil {
		log.Println(err)

		return "", err
	}
	defer f.Close()

	var errEncodingImage error
	switch format {
	case "jpeg":
		errEncodingImage = jpeg.Encode(f, img, nil)
	case "png":
		errEncodingImage = png.Encode(f, img)
	case "gif":
		errEncodingImage = gif.Encode(f, img, nil)
	default:
		errEncodingImage = jpeg.Encode(f, img, nil)
	}

	if errEncodingImage != nil {
		log.Println(err)

		return "", errEncodingImage
	}

	if err != nil {
		log.Println(err)

		return "", err
	}

	return fmt.Sprintf("http://localhost:8080/FileAzure/Data/Image/%s/profile/%s", userId, filename), nil
}

// Updating the users image and delete the old image
func UpdateUserProfileImageBool(userId string, newImgUrl string, oldImgUrl string) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE user_profile SET image_url = $1 WHERE user_id = $2`

	_, err := db.Exec(sqlStatement, newImgUrl, userId)

	if err != nil {
		log.Fatalf("\nUPDATE USER PROFILE IMAGE - Cannot execute command : %v\n", err)
		return false, err
	}

	DelOldImgUrl, errDelOldImgUrl := DeleteUsersFileImage(oldImgUrl)

	if errDelOldImgUrl != nil {
		return false, errDelOldImgUrl
	}

	if !DelOldImgUrl {
		return false, fmt.Errorf("%s", "Cannot Delete the old img, contact dev!")
	}

	return true, nil
}

// Delete users image url from database
func DeleteUserImageProfileBool(userId, oldImageUrl string) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE user_profile SET image_url = '' WHERE user_id = $1`

	_, err := db.Exec(sqlStatement, userId)

	if err != nil {
		log.Fatalf("\nDELETE USER PROFILE IMAGE - Cannot execute command : %v\n", err)
		return false, err
	}

	DelOldImgUrl, errDelOldImgUrl := DeleteUsersFileImage(oldImageUrl)

	if errDelOldImgUrl != nil {
		return false, errDelOldImgUrl
	}

	if !DelOldImgUrl {
		return false, fmt.Errorf("%s", "Cannot Delete the old img, contact dev!")
	}

	return true, nil
}

/*
Note : You should store the file path using env variable to secure the path
Example *export SECRET_PATH=/path/to/your/secret/path*

then use this to call the path : filePath := os.Getenv("SECRET_FILE_PATH")
*/
// Remove the users old image
func DeleteUsersFileImage(oldImageUrl string) (bool, error) {

	u, err := url.Parse(oldImageUrl)
	if err != nil {
		return false, err
	}
	path := u.Path

	// Extract the folder and filename from the path
	folders := strings.Split(path, "/")
	if len(folders) < 3 {
		return false, fmt.Errorf("%s %s", "Invalid image path: ", oldImageUrl)
	}
	folder := folders[4]
	filename := folders[len(folders)-1]

	// Delete the file
	err = os.Remove(fmt.Sprintf("/AzureWS/FileAzure/Data/Image/%s/profile/%s", folder, filename))
	if err != nil {
		return false, err
	}

	return true, nil
}
