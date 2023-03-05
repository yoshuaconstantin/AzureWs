package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	Aunth "AzureWS/globalvariable/authenticator"
	"AzureWS/globalvariable/constant"
	"AzureWS/logging"
	"AzureWS/module"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"

)

func GetAllCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query()

	tokenParam := queryParams.Get("token")
	indexParam := queryParams.Get("index")

	index, errI := strconv.Atoi(indexParam)

	if errI != nil {

		logging.InsertLog(r,constant.CommunityPost, errI.Error(), "", http.StatusBadRequest, 1, 1)

		http.Error(w, errI.Error(), http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, tokenParam)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPost, errAunth.Error(), tokenParam, http.StatusBadRequest, 1, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	getAllPostData, errGetDat := module.GetAllCommunityPostFromDB(index, GetUserIdAunth)

	if errGetDat != nil {

		logging.InsertLog(r,constant.CommunityPost, errGetDat.Error(), tokenParam, http.StatusUnauthorized, 1, 3)

		http.Error(w, errGetDat.Error(), http.StatusUnauthorized)
		return
	}

	logging.InsertLog(r,constant.CommunityPost, "", tokenParam, http.StatusOK, 1, 4)

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

		logging.InsertLog(r,constant.CommunityPost, err.Error(), "", http.StatusBadRequest, 2, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, postData.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPost, errAunth.Error(), postData.Token, AunthStatus, 2, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	var usrPrflMdl models.GetUserProfileDataModel

	GetDataProfile, errGetDatP := module.GetUserProfileDataFromDB(GetUserIdAunth)

	if errGetDatP != nil {

		logging.InsertLog(r,constant.CommunityPost, errGetDatP.Error(), postData.Token, AunthStatus, 2, 3)

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

		logging.InsertLog(r,constant.CommunityPost, errInsert.Error(), postData.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, errInsert.Error(), http.StatusInternalServerError)
		return
	}

	if !insrtNewCommunityPost {

		logging.InsertLog(r,constant.CommunityPost, "Cannot Insert New Post!", postData.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, "Cannot Insert New Post!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes creating new post"

	logging.InsertLog(r,constant.CommunityPost, "", postData.Token, http.StatusOK, 2, 4)


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestUpdateP request.RequestUpdateCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestUpdateP)
	if err != nil {

		logging.InsertLog(r,constant.CommunityPost, err.Error(), "", http.StatusOK, 3, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestUpdateP.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPost, errAunth.Error(), requestUpdateP.Token, AunthStatus, 3, 2)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	updateUsersCommunityPost, errUpdate := module.UpdateCommunityPostFromDB(GetUserIdAunth, requestUpdateP.Data[0])

	if errUpdate != nil {

		logging.InsertLog(r,constant.CommunityPost, errUpdate.Error(), requestUpdateP.Token, http.StatusInternalServerError, 1, 3)

		http.Error(w, errUpdate.Error(), http.StatusInternalServerError)
		return
	}

	if !updateUsersCommunityPost {

		logging.InsertLog(r,constant.CommunityPost, "Cannot Update Post!", requestUpdateP.Token, AunthStatus, 3, 3)

		http.Error(w, "Cannot Update Post, try again later!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes updating post"

	logging.InsertLog(r,constant.CommunityPost, "", requestUpdateP.Token, http.StatusInternalServerError, 3, 4)


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUserCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestDeletP request.RequestDeleteCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestDeletP)
	if err != nil {

		logging.InsertLog(r,constant.CommunityPost, err.Error(), "", http.StatusInternalServerError, 4, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestDeletP.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPost, errAunth.Error(), requestDeletP.Token, AunthStatus, 4, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	deleteUsersCommunityPost, errDeletP := module.DeleteCommunityPostFromDB(GetUserIdAunth, requestDeletP.Data[0])

	if errDeletP != nil {

		logging.InsertLog(r,constant.CommunityPost, errDeletP.Error(), requestDeletP.Token, http.StatusInternalServerError, 4, 3)

		http.Error(w, errDeletP.Error(), http.StatusInternalServerError)
		return
	}

	if !deleteUsersCommunityPost {

		logging.InsertLog(r,constant.CommunityPost, "Cannot Delete Post!", requestDeletP.Token, http.StatusInternalServerError, 4, 3)

		http.Error(w, "Cannot Delete Post, try again later!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes deleting post"

	logging.InsertLog(r,constant.CommunityPost, "", requestDeletP.Token, http.StatusOK, 4, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetSpecificCommunityPostComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestData request.RequestGetCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {

		logging.InsertLog(r,constant.CommunityPostCommentGet, err.Error(), "", http.StatusOK, 2, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestData.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPostCommentGet, errAunth.Error(), requestData.Token, AunthStatus, 2, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	getPostCmntDat, errGetDat := module.GetSpecificCommentCommunityPostFromDB(requestData.Data[0].PostId, requestData.Data[0].Index, GetUserIdAunth)

	if errGetDat != nil {

		logging.InsertLog(r,constant.CommunityPostCommentGet, errGetDat.Error(), requestData.Token, http.StatusConflict, 2, 3)

		http.Error(w, errGetDat.Error(), http.StatusConflict)
		return
	}

	var response response.ResponseGetCommunityPostComment
	response.Status = http.StatusOK
	response.Message = "Succes"
	response.Data = getPostCmntDat

	logging.InsertLog(r,constant.CommunityPostCommentGet, "", requestData.Token, http.StatusOK, 2, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func InsertNewCommentCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestInsrtCmnt request.RequestInsertCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestInsrtCmnt)
	if err != nil {

		logging.InsertLog(r,constant.CommunityPostComment, err.Error(), "", http.StatusBadRequest, 2, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestInsrtCmnt.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPostComment, errAunth.Error(), requestInsrtCmnt.Token, AunthStatus, 2, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	insrtNewCmntCmmntyP, errInsertCmt := module.InsertCommentCommunityPostToDB(GetUserIdAunth, requestInsrtCmnt.Data[0])

	if errInsertCmt != nil {

		logging.InsertLog(r,constant.CommunityPostComment, errInsertCmt.Error(), requestInsrtCmnt.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, errInsertCmt.Error(), http.StatusInternalServerError)
		return
	}

	if !insrtNewCmntCmmntyP {

		logging.InsertLog(r,constant.CommunityPostComment, "Cannot Insert New Comment!", requestInsrtCmnt.Token, http.StatusInternalServerError, 2, 3)

		http.Error(w, "Cannot Insert New Comment", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes insert comment"

	logging.InsertLog(r,constant.CommunityPostComment, "", requestInsrtCmnt.Token, http.StatusOK, 2, 3)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserCommentCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestUpdateCmntP request.RequestUpdateCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestUpdateCmntP)
	if err != nil {

		logging.InsertLog(r,constant.CommunityPostComment, err.Error(), "", http.StatusBadRequest, 3, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestUpdateCmntP.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPostComment, errAunth.Error(), requestUpdateCmntP.Token, AunthStatus, 3, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	updateUsersCmntCmntyPost, errUpdateCmnt := module.UpdateCommentCommunityPostFromDB(GetUserIdAunth, requestUpdateCmntP.Data[0])

	if errUpdateCmnt != nil {

		logging.InsertLog(r,constant.CommunityPostComment, errUpdateCmnt.Error(), requestUpdateCmntP.Token, http.StatusInternalServerError, 3, 3)

		http.Error(w, errUpdateCmnt.Error(), http.StatusInternalServerError)
		return
	}

	if !updateUsersCmntCmntyPost {

		logging.InsertLog(r,constant.CommunityPostComment, errUpdateCmnt.Error(), requestUpdateCmntP.Token, http.StatusInternalServerError, 3, 3)


		http.Error(w, "Cannot Update User Comment!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes updating comment"

	logging.InsertLog(r,constant.CommunityPostComment, "", requestUpdateCmntP.Token, http.StatusOK, 3, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUserCommentCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestDeletCmntP request.RequestDeleteCommentCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestDeletCmntP)
	if err != nil {

		logging.InsertLog(r,constant.CommunityPostComment, err.Error(), "", http.StatusOK, 4, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestDeletCmntP.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPostComment, errAunth.Error(), requestDeletCmntP.Token, AunthStatus, 4, 3)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	deleteUsersCmntCmntyPost, errDeletCmntP := module.DeleteCommentCommunityPostFromDB(GetUserIdAunth, requestDeletCmntP.Data[0])

	if errDeletCmntP != nil {

		logging.InsertLog(r,constant.CommunityPostComment, errDeletCmntP.Error(), requestDeletCmntP.Token,  http.StatusInternalServerError, 4, 3)

		http.Error(w, errDeletCmntP.Error(), http.StatusInternalServerError)
		return
	}

	if !deleteUsersCmntCmntyPost {

		logging.InsertLog(r,constant.CommunityPostComment, "cannot Delete Post!", requestDeletCmntP.Token,  http.StatusInternalServerError, 4, 3)

		http.Error(w, "Cannot Delete Post, try again later!", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes deleting comment"

	logging.InsertLog(r,constant.CommunityPostComment, "", requestDeletCmntP.Token,  http.StatusOK, 4, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func InsertNewLikeCommunityPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestLike request.RequestInsertLikeCommunityPost

	err := json.NewDecoder(r.Body).Decode(&requestLike)
	if err != nil {

		logging.InsertLog(r,constant.CommunityPostLike, err.Error(), "",  http.StatusBadRequest, 2, 2)

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	GetUserIdAunth, AunthStatus, errAunth := Aunth.SecureAuthenticator(w, r, requestLike.Token)

	if errAunth != nil {

		logging.InsertLog(r,constant.CommunityPostLike, errAunth.Error(), requestLike.Token,  AunthStatus, 2, 2)

		http.Error(w, errAunth.Error(), AunthStatus)
		return
	}

	likeCmntyP, errLiktCmntyP := module.InsertLikeCommunityPostToDB(GetUserIdAunth, requestLike.Data[0])

	if errLiktCmntyP != nil {

		logging.InsertLog(r,constant.CommunityPostLike, errLiktCmntyP.Error(), requestLike.Token,   http.StatusInternalServerError, 2, 3)

		http.Error(w, errLiktCmntyP.Error(), http.StatusInternalServerError)
		return
	}

	if !likeCmntyP {

		logging.InsertLog(r,constant.CommunityPostLike,"cannot Like Post", requestLike.Token,   http.StatusInternalServerError, 2, 3)
		
		http.Error(w, "Cannot Like post", http.StatusInternalServerError)
		return
	}

	var response response.GeneralResponseNoData
	response.Status = http.StatusOK
	response.Message = "Succes"

	logging.InsertLog(r,constant.CommunityPostLike,"", requestLike.Token,   http.StatusOK, 2, 4)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
