package controller

import (
	"encoding/json" // package untuk enkode dan mendekode json menjadi struct dan sebaliknya
	"log"
	"net/http" // digunakan untuk mengakses objek permintaan dan respons dari api

	"github.com/gorilla/mux" // digunakan untuk mendapatkan parameter dari router
	_ "github.com/lib/pq"    // postgres golang driver

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

type GetDashboards struct {
	Token string `json:"token"`
}

func GetDshbrdDat(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	token := params["token"]

	userId, errGetUuid := validation.ValidateTokenGetUuid(token)

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

func UpdtDshbrdDat(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	token := params["token"]
	mode := params["mode"]

	boolResult, err := models.UpdateDashboardsData(token, mode)

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
