package controller

import (
	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/module"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

// Get User Dashboards Data based on current user
func GetDashboardsData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	datas, err := module.GetDashboardsDataFromDB(GetUserIdAunth)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	var response response.ResponseDashboardsData
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = datas

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update User Dashboards Data base on current user
func UpdateDashboardsData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var modelUpdate request.RequestUpdateDashboardsData
	err := json.NewDecoder(r.Body).Decode(&modelUpdate)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, modelUpdate.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	boolResult, err := module.UpdateDashboardsDataFromDB(GetUserIdAunth, modelUpdate.Mode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !boolResult {
		http.Error(w, "Cannot update data", http.StatusInternalServerError)
		return
	} 

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
