package controller

import (
	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/module"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"encoding/json"
	"net/http"
)

/*
<Documentation>
	Step-by-step how to use this (currently untested yet)
	- Upload an img byte from apps to UploadImage endpoint -> if success then server will return response with ImageUrl to user
	- Hit InsertDataProfile endpoint and store the remaining data with stored ImageUrl from Upload
	- Leave / Refresh the apps then hit GetDataProfile to get the whole data
** Alt.Step : Hit only UploadImage to get image only to database (untested, unchecked flow)
</Documentation>
*/

// Upload image to local storage outside of source code with return to user ImageUrl string
func UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Read the binary data from the request body
	var imageData request.RequestUploadImageData
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, imageData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	// Func Convert from byte to image and return string.
	UploadImageToDB, errUploadImageToDB := module.UploadUserProfilePhotoToDB(GetUserIdAunth, imageData.Data)

	if errUploadImageToDB != nil {
		http.Error(w, errUploadImageToDB.Error(), http.StatusNotAcceptable)
		return
	}

	if !UploadImageToDB {
		http.Error(w, "Failed to upload Image Url to database", http.StatusNotAcceptable)
		return
	}

	// Send back the response to user with file name
	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success Upload Image"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update image profile, replace the string
func UpdateImageProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Read the binary data from the request body
	var imageData request.RequestUpdateProfileImageData
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, imageData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	// Func Convert from byte to image and return string.
	GetNewImageUrl, errGetImageUrl := module.ConvertByteToImgString(imageData.Data, GetUserIdAunth)

	if errGetImageUrl != nil {
		http.Error(w, errGetImageUrl.Error(), http.StatusNotAcceptable)
		return
	}

	UpdateUsersProfileImage, errUpdateProfileImage := module.UpdateUserProfileImageFromDB(GetUserIdAunth, GetNewImageUrl, imageData.OldImageUrl)

	if errUpdateProfileImage != nil {
		http.Error(w, errUpdateProfileImage.Error(), http.StatusNotAcceptable)
		return
	}

	if !UpdateUsersProfileImage {
		http.Error(w, "An error occured when updating users profile image", http.StatusNotAcceptable)
		return
	}

	// Send back the response to user with file name
	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success Update Image"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete image profile, replace with empty string
func DeleteImageProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Read the binary data from the request body
	var imageData request.RequestDeleteProfileImageData
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, imageData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	DeleteUserProfileImage, errDeleteProfileImage := module.DeleteUserImageProfileFromDB(GetUserIdAunth, imageData.OldImageUrl)

	if errDeleteProfileImage != nil {
		http.Error(w, errDeleteProfileImage.Error(), http.StatusNotAcceptable)
		return
	}

	if !DeleteUserProfileImage {
		http.Error(w, "An error occured when deleting users profile image", http.StatusNotAcceptable)
		return
	}

	// Send back the response to user with file name
	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success Update Image"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Insert profile data into database
func InsertDataProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var userData request.RequestInsertProfileData

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, userData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	data := userData.Data[0]

	// Create a UserProfileData instance using the data from the InsertProfileData struct
	userProfileData := models.UserProfileDataModel{
		Nickname: &data.Nickname,
		Age:      &data.Age,
		Gender:   &data.Gender,
		ImageUrl: &data.ImageUrl,
	}

	insertDataProfile, errInsertUserProfile := module.UpdateUserProfileToDatabase(userProfileData, GetUserIdAunth)

	if errInsertUserProfile != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusConflict
		response.Message = errInsertUserProfile.Error()

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)

		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = insertDataProfile

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get the profile data with path to local dir
func GetDataProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var TokenData request.RequestTokenOnlyData
	err := json.NewDecoder(r.Body).Decode(&TokenData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, TokenData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	var dataModel []models.GetUserProfileDataModel

	dataModel, errGetDataProfile := module.GetUserProfileDataFromDB(GetUserIdAunth)

	if errGetDataProfile != nil {
		http.Error(w, errGetDataProfile.Error(), http.StatusConflict)
		return
	}

	var response response.ResponseGetDataUserProfile
	response.Status = http.StatusOK
	response.Message = "Get data success"
	response.Data = dataModel

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
