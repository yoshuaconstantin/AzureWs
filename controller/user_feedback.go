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

	_ "github.com/lib/pq" // postgres golang driver
)

// Insert Feedback User func
func InsrtFdbckUsr(w http.ResponseWriter, r *http.Request) {
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

	var GetDataRequest request.RequestInsertFeedback
	err := json.NewDecoder(r.Body).Decode(&GetDataRequest)
	if err != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenEntered := GetDataRequest.Token

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

	GetNickname, errGetNick := validation.ValidateGetNickname(userId)

	if errGetNick != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = errGetNick.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	InsertUserFeedback, errInsertFeedback := module.InsertFeedbackUserToDB(userId, GetNickname, GetDataRequest.Comment, GetDataRequest.Timestamp)

	if errInsertFeedback != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = errInsertFeedback.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !InsertUserFeedback {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Failed to insert feedback"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update Comment Feedback User func
func UpdtCmmtFdbckUsr(w http.ResponseWriter, r *http.Request) {
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

	var GetDataRequest request.RequestEditFeedback
	err := json.NewDecoder(r.Body).Decode(&GetDataRequest)
	if err != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenEntered := GetDataRequest.Token

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

	EditFeedback, errEditFeedback := module.EditFeedBackUserFromDB(GetDataRequest.Id, GetDataRequest.Comment, GetDataRequest.Timestamp, userId)

	if errEditFeedback != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = errEditFeedback.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !EditFeedback {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Failed to Edit feedback"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get All feedback data
func GetFdbckUsr(w http.ResponseWriter, r *http.Request) {
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

	var GetDataRequest request.RequestGetAllFeedbackData
	err := json.NewDecoder(r.Body).Decode(&GetDataRequest)
	if err != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenEntered := GetDataRequest.Token

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

	var offset = GetDataRequest.Index * 10

	GetAllFeedbackData, errGetFeedbackData := module.GetFeedBackUserDataFromDB(offset)

	if errGetFeedbackData != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = errGetFeedbackData.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response response.ResponseGetAllFeedBackUser
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = GetAllFeedbackData

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete Users Feedback func
func DelUsrFdbck(w http.ResponseWriter, r *http.Request) {
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

	var GetDataRequest request.RequestDeleteSingleFeedbackData
	err := json.NewDecoder(r.Body).Decode(&GetDataRequest)
	if err != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusBadRequest
		response.Message = "Invalid request body"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenEntered := GetDataRequest.Token

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

	DeleteSingleFeedBackUser, errDeleteSinglFdbck := module.DeleteFeedBackUserFromDB(GetDataRequest.Id, userId)

	if errDeleteSinglFdbck != nil {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = errDeleteSinglFdbck.Error()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !DeleteSingleFeedBackUser {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Cannot delete this feedback"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
