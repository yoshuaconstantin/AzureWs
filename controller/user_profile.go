package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"AzureWS/models"
)

type responseUserProfile struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Read the binary data from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Decode the binary data as an image
	img, format, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding image", http.StatusBadRequest)
		return
	}

	// Choose a unique filename for the image
	filename := fmt.Sprintf("%d.%s", time.Now().Unix(), format)
	filepath := filepath.Join("images", filename)

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

	// Return success
	fmt.Fprintf(w, "UPLOAD IMAGE - Image uploaded successfully: %s", filename)

	var response responseUserProfile
	response.Status = http.StatusOK
	response.Message = "Succes Upload Image"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

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

	insertDataProfile, errInsertUserProfile := models.InsertUserProfileToDatabase(userData)

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
