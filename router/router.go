package router

import (
	"AzureWS/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	// User_Login API
	router.HandleFunc("/api/users", controller.GetAllUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{id}", controller.GetSnglUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/add_user", controller.InsrtNewUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user", controller.UpdtUserPsswd).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/user", controller.DltUsr).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/login", controller.LoginUser).Methods("GET", "OPTIONS")

	// Dashboards data using token
	router.HandleFunc("/api/home/dashboards", controller.GetDshbrdDat).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/home/update/dashboard/data", controller.UpdtDshbrdDat).Methods("POST", "OPTIONS")

	// User Profile using token
	router.HandleFunc("/api/home/user/profile/image", controller.UploadImage).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/profile/image", controller.UpdateImageProfile).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/home/user/profile/image", controller.DeleteImageProfile).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/home/user/profile", controller.InsertDataProfile).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/profile", controller.GetDataProfile).Methods("GET", "OPTIONS")

	return router
}
