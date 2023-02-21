package controller

import (
	jwttoken "AzureWS/JWTTOKEN"
	"AzureWS/module"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"AzureWS/session"
	"AzureWS/validation"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// Get User Dashboards Data based on current user
func GetDshbrdDat(w http.ResponseWriter, r *http.Request) {
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

	var token request.RequestToken
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		var response response.GeneralResponseNoData
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

		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving UserId"

		w.WriteHeader(http.StatusInternalServerError)
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

	SessionValidation, errSessionCheck := session.CheckSessionInside(userId)

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

	datas, err := module.GetDashboardsDataFromDB(userId)

	if err != nil {
		log.Fatalf("Unable to retrieve data. %v", err)

		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving data"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)

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
func UpdtDshbrdDat(w http.ResponseWriter, r *http.Request) {

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

	var modelUpdate request.RequestUpdateDashboardsData
	err := json.NewDecoder(r.Body).Decode(&modelUpdate)
	if err != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	userId, errGetUuid := validation.ValidateTokenGetUuid(modelUpdate.Token)

	if errGetUuid != nil {
		log.Fatalf("Unable to retrieve UserId. %v", errGetUuid)

		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving UserId"

		w.WriteHeader(http.StatusInternalServerError)
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

	SessionValidation, errSessionCheck := session.CheckSessionInside(userId)

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

	boolResult, err := module.UpdateDashboardsDataFromDB(userId, modelUpdate.Mode)

	if err != nil {
		log.Fatalf("Unable to retrieve data. %v", err)

		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Error retrieving data"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if boolResult {
		var response response.GeneralResponseNoData
		response.Status = http.StatusOK
		response.Message = "Success"

		// Send the response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else {
		var response response.GeneralResponseNoData
		response.Status = http.StatusNotAcceptable
		response.Message = "Error try again later / contact dev"

		// Send the response
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}
}
