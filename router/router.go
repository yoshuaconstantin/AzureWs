package router

import (
	"github.com/gorilla/mux"

	"AzureWS/controller"
	"AzureWS/websocketstruct"

)

func Router() *mux.Router {

	router := mux.NewRouter()

	// User_Login API
	// router.HandleFunc("/api/users", controller.GetAllUsr).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/user/{id}", controller.GetSnglUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/create-account", controller.CreateNewAccount).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user", controller.UpdateAccountPassword).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/user", controller.DeleteAccount).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/login", controller.LoginAccount).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/logout", controller.LogoutAccount).Methods("GET", "OPTIONS")

	// Validation data
	router.HandleFunc("/api/password-validation", controller.ValidatePassword).Methods("GET", "OPTIONS")

	// Dashboards data using token and JWT
	router.HandleFunc("/api/home/dashboards", controller.GetDashboardsData).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/home/update/dashboard/data", controller.UpdateDashboardsData).Methods("POST", "OPTIONS")

	// User Profile using token and JWT
	router.HandleFunc("/api/home/user/profile/image", controller.UploadImage).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/profile/image", controller.UpdateImageProfile).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/home/user/profile/image", controller.DeleteImageProfile).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/home/user/profile", controller.InsertDataProfile).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/profile", controller.GetDataProfile).Methods("GET", "OPTIONS")

	// JWT Stuff
	router.HandleFunc("/api/token-refresh", controller.RefrshToken).Methods("GET", "OPTIONS")

	// Testing JWT
	router.HandleFunc("/api/generate", controller.TestGenerateJwt).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/verify", controller.TestVerifyJwt).Methods("GET", "OPTIONS")

	// Feedback user
	router.HandleFunc("/api/home/user/feedback", controller.GetAllFeedbackUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/home/user/feedback", controller.InsertFeedbackUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/feedback", controller.UpdateCommentFeedbackUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/home/user/feedback", controller.DeletUserFeedback).Methods("DELETE", "OPTIONS")

	// Community chat using websocket
	router.HandleFunc("/community_chat", websocketstruct.CommunityChat)

	// Community Post
	router.HandleFunc("/api/community/post", controller.GetAllCommunityPost).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/community/post", controller.InsertNewCommunityPost).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/community/post", controller.UpdateUserCommunityPost).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/community/post", controller.DeleteUserCommunityPost).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/community/post/like", controller.InsertNewLikeCommunityPost).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/community/post/comment", controller.GetSpecificCommunityPostComment).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/community/post/comment", controller.InsertNewCommentCommunityPost).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/community/post/comment", controller.UpdateUserCommentCommunityPost).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/community/post/comment", controller.DeleteUserCommentCommunityPost).Methods("DELETE", "OPTIONS")
	return router

}
