package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"AzureWS/models"
	"AzureWS/session"
	"AzureWS/validation"
)

type responseUserProfile struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

type responseUserProfileData struct {
	Message string                      `json:"message,omitempty"`
	Status  int                         `json:"status,omitempty"`
	Data    []models.GetUserProfileData `json:"data,omitempty"`
}

type responseUserProfileUploadImage struct {
	Message  string `json:"message,omitempty"`
	Status   int    `json:"status,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
}

type uploadImageData struct {
	Token string `json:"token"`
	Data  []byte `json:"data"`
}

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
	var imageData uploadImageData
	if err := json.NewDecoder(r.Body).Decode(&imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	getUserId, errTokenValidate := validation.ValidateTokenGetUuid(imageData.Token)

	if errTokenValidate != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	sessionValidation, errSessionValidate := session.CheckSessionInside(getUserId)

	if errSessionValidate != nil {
		http.Error(w, errSessionValidate.Error(), http.StatusForbidden)
	}

	if !sessionValidation {
		http.Error(w, "Session Expired", http.StatusForbidden)
	}

	img, format, err := image.Decode(bytes.NewReader(imageData.Data))
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding image", http.StatusBadRequest)
		return
	}

	path := filepath.Join("FileAzure", "Data", "Image", getUserId, "profile")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Choose a unique filename for the image
	filename := fmt.Sprintf("%d.%s", time.Now().Unix(), format)
	filepath := filepath.Join(path, filename)

	// Create the file and write the image to it
	f, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		http.Error(w, "UPLOAD IMAGE - Error creating file", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	if format == "jpeg" {
		err = jpeg.Encode(f, img, nil)
	} else {
		http.Error(w, "UPLOAD IMAGE - Format invalid", http.StatusInternalServerError)
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "UPLOAD IMAGE - Error encoding image", http.StatusInternalServerError)
		return
	}

	// Send back the response to user with file name
	var response responseUserProfileUploadImage
	response.Status = http.StatusOK
	response.Message = "Success Upload Image"
	response.ImageUrl = fmt.Sprintf("http://localhost:8080/FileAzure/Data/Image/%s/profile/%s", getUserId, filename)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Insert profile data into database
func InsertDataProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var userData models.UserProfileData

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		var response responseUserProfile
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	GetUserID, errGetUuid := validation.ValidateTokenGetUuid(*userData.Token)

	if errGetUuid != nil {
		var response responseUserProfile
		response.Status = http.StatusUnauthorized
		response.Message = errGetUuid.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	}

	SessionValidation, errSessionCheck := session.CheckSessionInside(GetUserID)

	if errSessionCheck != nil {
		var response responseUserProfile
		response.Status = http.StatusForbidden
		response.Message = errSessionCheck.Error()

		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
	}

	if !SessionValidation {
		var response responseUserProfile
		response.Status = http.StatusUnauthorized
		response.Message = "Session Expired"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	}

	insertDataProfile, errInsertUserProfile := models.InsertUserProfileToDatabase(userData, GetUserID)

	if errInsertUserProfile != nil {
		var response responseUserProfile
		response.Status = http.StatusConflict
		response.Message = errInsertUserProfile.Error()

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response responseUserProfile
	response.Status = http.StatusOK
	response.Message = insertDataProfile

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get the profile data with path to local dir
func GetDataProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var token string
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		var response responseUserProfile
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	getUserId, errGetUuid := validation.ValidateTokenGetUuid(token)

	if errGetUuid != nil {
		var response responseUserProfile
		response.Status = http.StatusUnauthorized
		response.Message = errGetUuid.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	}

	checkSession, errCheckSession := session.CheckSessionInside(getUserId)

	if errCheckSession != nil {
		var response responseUserProfile
		response.Status = http.StatusNotFound
		response.Message = errCheckSession.Error()

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
	}

	if checkSession {
		var dataModel []models.GetUserProfileData

		dataModel, errGetDataProfile := models.GetUserProfileDataFromDatabase(getUserId)

		if errGetDataProfile != nil {
			var response responseUserProfile
			response.Status = http.StatusConflict
			response.Message = errGetDataProfile.Error()

			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
		}

		var response responseUserProfileData
		response.Status = http.StatusOK
		response.Message = "Get data success"
		response.Data = dataModel

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	} else {
		var response responseUserProfile
		response.Status = http.StatusUnauthorized
		response.Message = "Session Expired"

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	}
}
