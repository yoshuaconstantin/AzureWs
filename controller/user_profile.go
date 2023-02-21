package controller

import (
	jwttoken "AzureWS/JWTTOKEN"
	"AzureWS/module"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"AzureWS/session"
	"AzureWS/validation"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// If the authorization header is empty, return an error
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Missing authorization header"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Read the binary data from the request body
	var imageData request.RequestUploadImageData
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getUserId, errTokenValidate := validation.ValidateTokenGetUuid(imageData.Token)

	if errTokenValidate != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(tokenString)

	if erroCheckJWt != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = erroCheckJWt.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckJwtTokenValidation {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Unauthorized user"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	sessionValidation, errSessionValidate := session.CheckSessionInside(getUserId)

	if errSessionValidate != nil {
		http.Error(w, errSessionValidate.Error(), http.StatusForbidden)
		return
	}

	if !sessionValidation {
		http.Error(w, "Session Expired", http.StatusForbidden)
		return
	}

	// Func Convert from byte to image and return string.
	UploadImageToDB, errUploadImageToDB := module.UploadUserProfilePhotoBool(getUserId, imageData.Data)

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

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// If the authorization header is empty, return an error
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Missing authorization header"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Read the binary data from the request body
	var imageData request.RequestUpdateProfileImageData
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getUserId, errTokenValidate := validation.ValidateTokenGetUuid(imageData.Token)

	if errTokenValidate != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(tokenString)

	if erroCheckJWt != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = erroCheckJWt.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckJwtTokenValidation {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Unauthorized user"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	sessionValidation, errSessionValidate := session.CheckSessionInside(getUserId)

	if errSessionValidate != nil {
		http.Error(w, errSessionValidate.Error(), http.StatusForbidden)
		return
	}

	if !sessionValidation {
		http.Error(w, "Session Expired", http.StatusForbidden)
		return
	}

	// Func Convert from byte to image and return string.
	GetNewImageUrl, errGetImageUrl := module.ConvertByteToImgString(imageData.Data, getUserId)

	if errGetImageUrl != nil {
		http.Error(w, errGetImageUrl.Error(), http.StatusNotAcceptable)
		return
	}

	UpdateUsersProfileImage, errUpdateProfileImage := module.UpdateUserProfileImageBool(getUserId, GetNewImageUrl, imageData.OldImageUrl)

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

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// If the authorization header is empty, return an error
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Missing authorization header"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Read the binary data from the request body
	var imageData request.RequestDeleteProfileImageData
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getUserId, errTokenValidate := validation.ValidateTokenGetUuid(imageData.Token)

	if errTokenValidate != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(tokenString)

	if erroCheckJWt != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = erroCheckJWt.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckJwtTokenValidation {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Unauthorized user"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	sessionValidation, errSessionValidate := session.CheckSessionInside(getUserId)

	if errSessionValidate != nil {
		http.Error(w, errSessionValidate.Error(), http.StatusForbidden)
		return
	}

	if !sessionValidation {
		http.Error(w, "Session Expired", http.StatusForbidden)
		return
	}

	DeleteUserProfileImage, errDeleteProfileImage := module.DeleteUserImageProfileBool(getUserId, imageData.OldImageUrl)

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

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// If the authorization header is empty, return an error
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Missing authorization header"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var userData request.RequestInsertProfileData

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println(userData)
	data := userData.Data[0]

	// Create a UserProfileData instance using the data from the InsertProfileData struct
	userProfileData := models.UserProfileDataModel{
		Nickname: &data.Nickname,
		Age:      &data.Age,
		Gender:   &data.Gender,
		ImageUrl: &data.ImageUrl,
	}

	GetUserID, errGetUuid := validation.ValidateTokenGetUuid(userData.Token)

	if errGetUuid != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = errGetUuid.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(tokenString)

	if erroCheckJWt != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = erroCheckJWt.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckJwtTokenValidation {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Unauthorized user"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	SessionValidation, errSessionCheck := session.CheckSessionInside(GetUserID)

	if errSessionCheck != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusForbidden
		response.Message = errSessionCheck.Error()

		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !SessionValidation {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Session Expired"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	insertDataProfile, errInsertUserProfile := module.UpdateUserProfileToDatabase(userProfileData, GetUserID)

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

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// If the authorization header is empty, return an error
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Missing authorization header"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var TokenData request.RequestTokenOnlyData
	err := json.NewDecoder(r.Body).Decode(&TokenData)
	if err != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	getUserId, errGetUuid := validation.ValidateTokenGetUuid(TokenData.Token)

	if errGetUuid != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = errGetUuid.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	CheckJwtTokenValidation, erroCheckJWt := jwttoken.VerifyToken(tokenString)

	if erroCheckJWt != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = erroCheckJWt.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckJwtTokenValidation {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Unauthorized user"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	checkSession, errCheckSession := session.CheckSessionInside(getUserId)

	if errCheckSession != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusNotFound
		response.Message = errCheckSession.Error()

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if checkSession {
		var dataModel []models.GetUserProfileDataModel

		dataModel, errGetDataProfile := module.GetUserProfileDataFromDatabase(getUserId)

		if errGetDataProfile != nil {
			var response response.GeneralResponseNoData
			response.Status = http.StatusConflict
			response.Message = errGetDataProfile.Error()

			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		var response response.ResponseGetDataUserProfile
		response.Status = http.StatusOK
		response.Message = "Get data success"
		response.Data = dataModel

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return

	} else {
		var response response.GeneralResponseNoData
		response.Status = http.StatusUnauthorized
		response.Message = "Session Expired"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
}
