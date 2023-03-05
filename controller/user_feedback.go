package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/lib/pq" // postgres golang driver

	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/globalvariable/constant"
	"AzureWS/logging"
	"AzureWS/module"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
)

// Insert Feedback User func
func InsertFeedbackUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var GetDataRequest request.RequestInsertFeedback
	err := json.NewDecoder(r.Body).Decode(&GetDataRequest)
	if err != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, err.Error(), "", http.StatusBadRequest, 2, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, GetDataRequest.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errAunth.Error(), GetDataRequest.Token, AunthStatus, 2, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	GetProfileData, errGetDat := module.GetUserProfileDataFromDB(GetUserIdAunth)

	if errGetDat != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errGetDat.Error(), GetDataRequest.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, errGetDat.Error(), http.StatusInternalServerError)
		return
	}

	InsertUserFeedback, errInsertFeedback := module.InsertFeedbackUserToDB(GetUserIdAunth, *GetProfileData[0].Nickname, GetDataRequest.Comment)

	if errInsertFeedback != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errInsertFeedback.Error(), GetDataRequest.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, errInsertFeedback.Error(), http.StatusInternalServerError)
		return
	}

	if !InsertUserFeedback {

		logging.InsertLog(r, constant.HomeUserFeedback, "Failed to insert feedback", GetDataRequest.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, "Failed to insert feedback", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	logging.InsertLog(r, constant.HomeUserFeedback, "", GetDataRequest.Token, http.StatusOK, 2, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Update Comment Feedback User func
func UpdateCommentFeedbackUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var GetDataRequest request.RequestEditFeedback
	err := json.NewDecoder(r.Body).Decode(&GetDataRequest)
	if err != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, err.Error(), "", http.StatusOK, 3, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, GetDataRequest.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errAunth.Error(), GetDataRequest.Token, AunthStatus, 3, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	EditFeedback, errEditFeedback := module.EditFeedBackUserFromDB(GetDataRequest.Id, GetDataRequest.Comment, GetUserIdAunth)

	if errEditFeedback != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errEditFeedback.Error(), GetDataRequest.Token, http.StatusInternalServerError, 3, 3)

		http.Error(w, errEditFeedback.Error(), http.StatusInternalServerError)
		return
	}

	if !EditFeedback {

		logging.InsertLog(r, constant.HomeUserFeedback, "Failed to Edit feedback", GetDataRequest.Token, http.StatusInternalServerError, 3, 3)

		http.Error(w, "Failed to Edit feedback", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	logging.InsertLog(r, constant.HomeUserFeedback, "", GetDataRequest.Token, http.StatusOK, 3, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get All feedback data
func GetAllFeedbackUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")
	indexParam := queryParams.Get("index")

	index, errI := strconv.Atoi(indexParam)

	if errI != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errI.Error(), "", http.StatusBadRequest, 1, 1)

		http.Error(w, errI.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errAunth.Error(), tokenParam, AunthStatus, 1, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	var offset = index * 10

	GetAllFeedbackData, errGetFeedbackData := module.GetFeedBackUserDataFromDB(GetUserIdAunth, offset)

	if errGetFeedbackData != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errGetFeedbackData.Error(), tokenParam, http.StatusInternalServerError, 1, 3)

		http.Error(w, errGetFeedbackData.Error(), http.StatusInternalServerError)
		return
	}

	var response response.ResponseGetAllFeedBackUser
	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = GetAllFeedbackData

	logging.InsertLog(r, constant.HomeUserFeedback, "", tokenParam, http.StatusOK, 1, 3)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete Users Feedback func
func DeletUserFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var GetDataRequest request.RequestDeleteSingleFeedbackData
	err := json.NewDecoder(r.Body).Decode(&GetDataRequest)
	if err != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, err.Error(), "", http.StatusBadRequest, 4, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, GetDataRequest.Token)

	if errAunth != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errAunth.Error(), GetDataRequest.Token, AunthStatus, 4, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	DeleteSingleFeedBackUser, errDeleteSinglFdbck := module.DeleteFeedBackUserFromDB(GetDataRequest.Id, GetUserIdAunth)

	if errDeleteSinglFdbck != nil {

		logging.InsertLog(r, constant.HomeUserFeedback, errDeleteSinglFdbck.Error(), GetDataRequest.Token, http.StatusInternalServerError, 4, 3)

		http.Error(w, errDeleteSinglFdbck.Error(), http.StatusInternalServerError)
		return
	}

	if !DeleteSingleFeedBackUser {

		logging.InsertLog(r, constant.HomeUserFeedback, "Cannot delete this feedback", GetDataRequest.Token, http.StatusInternalServerError, 4, 3)

		http.Error(w, "Cannot delete this feedback", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Success"

	logging.InsertLog(r, constant.HomeUserFeedback, "", GetDataRequest.Token, http.StatusOK, 4, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
