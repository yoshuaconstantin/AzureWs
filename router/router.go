package router

import (
	"github.com/gorilla/mux"

	"AzureWS/controller"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	// User_Login API
	router.HandleFunc("/api/users", controller.GetAllUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{id}", controller.GetSnglUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/add_user", controller.InsrtNewUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", controller.UpdtUserPsswd).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/user/{id}", controller.DltUsr).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/login", controller.LoginUser).Methods("GET", "OPTIONS")

	// Dashboards data using token
	router.HandleFunc("/api/home/dashboards", controller.GetDshbrdDat).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/home/update/dashboard/data", controller.UpdtDshbrdDat).Methods("POST", "OPTIONS")

	// User Profile using token
	router.HandleFunc("/api/home/user/profile", controller.UploadImage).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/home/user/upload/image", controller.InsertDataProfile).Methods("POST", "OPTIONS")

	return router
}
