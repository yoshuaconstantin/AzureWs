package controller

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"

	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/globalvariable/constant"
	"AzureWS/logging"
	"AzureWS/module"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"

)

// Get User Dashboards Data based on current user
func GetDashboardsData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeDashboards, errAunth.Error(), "", AunthStatus, 1, 1)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	datas, err := module.GetDashboardsDataFromDB(GetUserIdAunth)

	if err != nil {

		logging.InsertLog(r, constant.HomeDashboards, err.Error(), tokenParam, http.StatusInternalServerError, 1, 3)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response response.ResponseDashboardsData
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = datas

	logging.InsertLog(r, constant.HomeDashboards, "", tokenParam, http.StatusOK, 1, 4)

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
	
		logging.InsertLog(r, constant.HomeDashboards, err.Error(), "", http.StatusBadRequest, 2, 2)
	
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, modelUpdate.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeDashboards, errAunth.Error(), "", AunthStatus, 2, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	boolResult, err := module.UpdateDashboardsDataFromDB(GetUserIdAunth, modelUpdate.Mode)

	if err != nil {

		logging.InsertLog(r, constant.HomeDashboards, err.Error(), modelUpdate.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !boolResult {

		logging.InsertLog(r, constant.HomeDashboards, "Cannot update data!", modelUpdate.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, "Cannot update data", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	logging.InsertLog(r, constant.HomeDashboards, "", modelUpdate.Token, http.StatusOK, 2, 4)

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
