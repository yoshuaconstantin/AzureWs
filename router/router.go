package router

import (
	"AzureWS/controller"
	"AzureWS/websocketstruct"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	// User_Login API
	// router.HandleFunc("/api/users", controller.GetAllUsr).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/user/{id}", controller.GetSnglUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/add_user", controller.InsrtNewUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user", controller.UpdtUserPsswd).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/user", controller.DltUsr).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/login", controller.LoginUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/logout", controller.LgoutUsr).Methods("GET", "OPTIONS")

	// Dashboards data using token and JWT
	router.HandleFunc("/api/home/dashboards", controller.GetDshbrdDat).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/home/update/dashboard/data", controller.UpdtDshbrdDat).Methods("POST", "OPTIONS")

	// User Profile using token and JWT
	router.HandleFunc("/api/home/user/profile/image", controller.UploadImage).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/profile/image", controller.UpdateImageProfile).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/home/user/profile/image", controller.DeleteImageProfile).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/home/user/profile", controller.InsertDataProfile).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/profile", controller.GetDataProfile).Methods("GET", "OPTIONS")

	// JWT Stuff
	router.HandleFunc("/api/token_refresh", controller.RefrshToken).Methods("GET", "OPTIONS")

	// Testing JWT
	router.HandleFunc("/api/generate", controller.TestGenerateJwt).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/verify", controller.TestVerifyJwt).Methods("GET", "OPTIONS")
	
	// Feedback user
	router.HandleFunc("/api/home/user/feedback", controller.GetFdbckUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/home/user/feedback", controller.InsrtFdbckUsr).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/feedback", controller.UpdtCmmtFdbckUsr).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/home/user/feedback", controller.DelUsrFdbck).Methods("DELETE", "OPTIONS")

	// Community chat using websocket
	router.HandleFunc("/community_chat", websocketstruct.CommunityChat)


	return router

}
