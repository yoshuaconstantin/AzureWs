package controller

import (
	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/module"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"encoding/json"
	"net/http"
)

func GetAllCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestData request.RequestTokenWithIndex

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	getAllPostData, errGetDat := module.GetAllCommunityPostFromDB(requestData.Index, GetUserIdAunth)

	if errGetDat != nil {
		http.Error(w, errGetDat.Error(), http.StatusUnauthorized)
		return
	}

	var response response.ResponseGetCommunityPost
	response.Status = http.StatusOK
	response.Message = "Succes get all post"
	response.Data = getAllPostData

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func InsertNewCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var postData request.RequestInsertCommunityPost

	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, postData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	var usrPrflMdl models.GetUserProfileDataModel

	GetDataProfile, errGetDatP := module.GetUserProfileDataFromDB(GetUserIdAunth)

	if errGetDatP != nil {
		http.Error(w, errGetDatP.Error(), http.StatusConflict)
		return
	}

	usrPrflMdl = GetDataProfile[0]

	// Create a UserProfileData instance using the data from the InsertProfileData struct
	communityPostData := models.PostDataModels{
		PostMessage: postData.PostMessage,
		Nickname:    *usrPrflMdl.Nickname,
		Nation:      *usrPrflMdl.Nation,
		ImageUrl:    *usrPrflMdl.ImageUrl,
	}

	insrtNewCommunityPost, errInsert := module.InsertCommunityPostToDB(GetUserIdAunth, communityPostData)

	if errInsert != nil {
		http.Error(w, errInsert.Error(), http.StatusInternalServerError)
		return
	}

	if !insrtNewCommunityPost {
		http.Error(w, "Cannot Insert New Post!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes creating new post"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestUpdateP request.RequestUpdateCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestUpdateP)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestUpdateP.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	updateUsersCommunityPost, errUpdate := module.UpdateCommunityPostFromDB(GetUserIdAunth, requestUpdateP.Data[0])

	if errUpdate != nil {
		http.Error(w, errUpdate.Error(), http.StatusInternalServerError)
		return
	}

	if !updateUsersCommunityPost {
		http.Error(w, "Cannot Update Post, try again later!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes updating post"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUserCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestDeletP request.RequestDeleteCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestDeletP)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestDeletP.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	deleteUsersCommunityPost, errDeletP := module.DeleteCommunityPostFromDB(GetUserIdAunth, requestDeletP.Data[0])

	if errDeletP != nil {
		http.Error(w, errDeletP.Error(), http.StatusInternalServerError)
		return
	}

	if !deleteUsersCommunityPost {
		http.Error(w, "Cannot Delete Post, try again later!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes deleting post"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetSpecificCommunityPostComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestData request.RequestGetCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestData.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	getPostCmntDat, errGetDat := module.GetSpecificCommentCommunityPostFromDB(requestData.Data[0].PostId, requestData.Data[0].Index, GetUserIdAunth)

	if errGetDat != nil {
		http.Error(w, errGetDat.Error(), http.StatusConflict)
		return
	}

	var response response.ResponseGetCommunityPostComment
	response.Status = http.StatusOK
	response.Message = "Succes"
	response.Data = getPostCmntDat

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func InsertNewCommentCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestInsrtCmnt request.RequestInsertCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestInsrtCmnt)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestInsrtCmnt.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	insrtNewCmntCmmntyP, errInsertCmt := module.InsertCommentCommunityPostToDB(GetUserIdAunth, requestInsrtCmnt.Data[0])

	if errInsertCmt != nil {
		http.Error(w, errInsertCmt.Error(), http.StatusInternalServerError)
		return
	}

	if !insrtNewCmntCmmntyP {
		http.Error(w, "Cannot Insert New Comment", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes insert comment"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserCommentCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestUpdateCmntP request.RequestUpdateCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestUpdateCmntP)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestUpdateCmntP.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	updateUsersCmntCmntyPost, errUpdateCmnt := module.UpdateCommentCommunityPostFromDB(GetUserIdAunth, requestUpdateCmntP.Data[0])

	if errUpdateCmnt != nil {
		http.Error(w, errUpdateCmnt.Error(), http.StatusInternalServerError)
		return
	}

	if !updateUsersCmntCmntyPost {
		var response response.GeneralResponseNoData
		response.Status = http.StatusInternalServerError
		response.Message = "Cannot Update Comment, try again later!"

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes updating comment"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUserCommentCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestDeletCmntP request.RequestDeleteCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestDeletCmntP)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestDeletCmntP.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	deleteUsersCmntCmntyPost, errDeletCmntP := module.DeleteCommentCommunityPostFromDB(GetUserIdAunth, requestDeletCmntP.Data[0])

	if errDeletCmntP != nil {
		http.Error(w, errDeletCmntP.Error(), http.StatusInternalServerError)
		return
	}

	if !deleteUsersCmntCmntyPost {
		http.Error(w, "Cannot Delete Post, try again later!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes deleting comment"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func InsertNewLikeCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestLike request.RequestInsertLikeCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestLike)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestLike.Token)

	if errAunth != nil {
		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	likeCmntyP, errLiktCmntyP := module.InsertLikeCommunityPostToDB(GetUserIdAunth, requestLike.Data[0])

	if errLiktCmntyP != nil {
		http.Error(w, errLiktCmntyP.Error(), http.StatusInternalServerError)
		return
	}

	if !likeCmntyP {
		http.Error(w, "Cannot Like post", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
