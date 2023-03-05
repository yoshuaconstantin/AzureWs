package controller

import (
	"encoding/json"
	"net/http"

	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/globalvariable/constant"
	"AzureWS/logging"
	"AzureWS/module"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"

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

		logging.InsertLog(r, constant.HomeUserProfileImage, err.Error(), "", http.StatusBadRequest, 2, 2)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, imageData.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserProfileImage, errAunth.Error(), imageData.Token, AunthStatus, 2, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	// Func Convert from byte to image and return string.
	UploadImageToDB, errUploadImageToDB := module.UploadUserProfilePhotoToDB(GetUserIdAunth, imageData.Data)

	if errUploadImageToDB != nil {

		logging.InsertLog(r, constant.HomeUserProfileImage, errUploadImageToDB.Error(), imageData.Token, http.StatusNotAcceptable, 2, 3)

		http.Error(w, errUploadImageToDB.Error(), http.StatusNotAcceptable)
		return
	}

	if !UploadImageToDB {

		logging.InsertLog(r, constant.HomeUserProfileImage, "Failed to upload Image Url to database", imageData.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, "Failed to upload Image Url to database", http.StatusInternalServerError)
		return
	}

	// Send back the response to user with file name
	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success Upload Image"

	logging.InsertLog(r, constant.HomeUserProfileImage, "", imageData.Token, http.StatusOK, 2, 4)

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

		logging.InsertLog(r, constant.HomeUserProfileImage, err.Error(), "", http.StatusBadRequest, 3, 2)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, imageData.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserProfileImage, errAunth.Error(), imageData.Token, AunthStatus, 3, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	// Func Convert from byte to image and return string.
	GetNewImageUrl, errGetImageUrl := module.ConvertByteToImgString(imageData.Data, GetUserIdAunth)

	if errGetImageUrl != nil {

		logging.InsertLog(r, constant.HomeUserProfileImage, errGetImageUrl.Error(), imageData.Token, http.StatusNotAcceptable, 3, 3)

		http.Error(w, errGetImageUrl.Error(), http.StatusNotAcceptable)
		return
	}

	UpdateUsersProfileImage, errUpdateProfileImage := module.UpdateUserProfileImageFromDB(GetUserIdAunth, GetNewImageUrl, imageData.OldImageUrl)

	if errUpdateProfileImage != nil {

		logging.InsertLog(r, constant.HomeUserProfileImage, errUpdateProfileImage.Error(), imageData.Token, http.StatusInternalServerError, 3, 3)

		http.Error(w, errUpdateProfileImage.Error(), http.StatusInternalServerError)
		return
	}

	if !UpdateUsersProfileImage {

		logging.InsertLog(r, constant.HomeUserProfileImage, "An error occured when updating users profile image", imageData.Token, http.StatusInternalServerError, 3, 3)

		http.Error(w, "An error occured when updating users profile image", http.StatusInternalServerError)
		return
	}

	// Send back the response to user with file name
	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success Update Image"

	logging.InsertLog(r, constant.HomeUserProfileImage, "", imageData.Token, http.StatusOK, 3, 4)

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

		logging.InsertLog(r, constant.HomeUserProfileImage, err.Error(), "", http.StatusBadRequest, 4, 2)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, imageData.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserProfileImage, errAunth.Error(), imageData.Token, AunthStatus, 4, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	DeleteUserProfileImage, errDeleteProfileImage := module.DeleteUserImageProfileFromDB(GetUserIdAunth, imageData.OldImageUrl)

	if errDeleteProfileImage != nil {

		logging.InsertLog(r, constant.HomeUserProfileImage, errDeleteProfileImage.Error(), imageData.Token, http.StatusInternalServerError, 4, 3)

		http.Error(w, errDeleteProfileImage.Error(), http.StatusInternalServerError)
		return
	}

	if !DeleteUserProfileImage {

		logging.InsertLog(r, constant.HomeUserProfileImage, "An error occured when deleting users profile image", imageData.Token, http.StatusInternalServerError, 4, 3)

		http.Error(w, "An error occured when deleting users profile image", http.StatusInternalServerError)
		return
	}

	// Send back the response to user with file name
	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success Update Image"

	logging.InsertLog(r, constant.HomeUserProfileImage, "", imageData.Token, http.StatusOK, 4, 4)

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

		logging.InsertLog(r, constant.HomeUserProfile, err.Error(), "", http.StatusBadRequest, 2, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, userData.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserProfile, errAunth.Error(),userData.Token, AunthStatus, 2, 3)

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

		logging.InsertLog(r, constant.HomeUserProfile, errInsertUserProfile.Error(), userData.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, errInsertUserProfile.Error(), http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = insertDataProfile

	logging.InsertLog(r, constant.HomeUserProfile, "", userData.Token, http.StatusOK, 2, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get the profile data with path to local dir
func GetDataProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserProfile, errAunth.Error(), tokenParam, AunthStatus, 1, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	var dataModel []models.GetUserProfileDataModel

	dataModel, errGetDataProfile := module.GetUserProfileDataFromDB(GetUserIdAunth)

	if errGetDataProfile != nil {

		logging.InsertLog(r, constant.HomeUserProfile, errGetDataProfile.Error(), tokenParam, http.StatusInternalServerError, 1, 3)

		http.Error(w, errGetDataProfile.Error(), http.StatusInternalServerError)
		return
	}

	var response response.ResponseGetDataUserProfile
	response.Status = http.StatusOK
	response.Message = "Get data success"
	response.Data = dataModel

	logging.InsertLog(r, constant.HomeUserProfile, "", tokenParam, http.StatusOK, 1, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
