package controller

import (
	"encoding/json" // package untuk enkode dan mendekode json menjadi struct dan sebaliknya
	"log"
	"net/http" // digunakan untuk mengakses objek permintaan dan respons dari api
	// digunakan untuk mendapatkan parameter dari router

	_ "github.com/lib/pq" // postgres golang driver

	"AzureWS/models" //models package dimana User didefinisikan
	"AzureWS/validation"

)

/*

 */

type responseDashboards struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

type ResponseDashboardsData struct {
	Status  int                     `json:"status"`
	Message string                  `json:"message"`
	Data    []models.DashboardsData `json:"data"`
}

type GetTokenUser struct {
	Token string `json:"token"`
}

type UpdateDashboardsData struct {
	Token string `json:"token"`
	Mode   string `json:"mode"`
}

// Get User Dashboards Data based on current user
func GetDshbrdDat(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var token GetTokenUser
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		var response responseDashboards
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenEntered := token.Token

	userId, errGetUuid := validation.ValidateTokenGetUuid(tokenEntered)

	if errGetUuid != nil {
		log.Fatalf("Unable to retrieve UserId. %v", errGetUuid)

		var response responseDashboards
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving UserId"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	datas, err := models.GetDashboardsData(userId)

	if err != nil {
		log.Fatalf("Unable to retrieve data. %v", err)

		var response responseDashboards
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving data"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	var response ResponseDashboardsData
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = datas

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update User Dashboards Data base on current user
func UpdtDshbrdDat(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var modelUpdate UpdateDashboardsData
	err := json.NewDecoder(r.Body).Decode(&modelUpdate)
	if err != nil {
		var response responseDashboards
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	userId, errGetUuid := validation.ValidateTokenGetUuid(modelUpdate.Token)

	if errGetUuid != nil {
		log.Fatalf("Unable to retrieve UserId. %v", errGetUuid)

		var response responseDashboards
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving UserId"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	boolResult, err := models.UpdateDashboardsData(userId, modelUpdate.Mode)

	if err != nil {
		log.Fatalf("Unable to retrieve data. %v", err)

		var response responseDashboards
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving data"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

		return
	}

	if boolResult {
		var response responseDashboards
		response.Status = http.StatusOK
		response.Message = "Success"

		// Send the response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		var response responseDashboards
		response.Status = http.StatusNotAcceptable
		response.Message = "Error try again later / contact dev"

		// Send the response
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
	}
}
