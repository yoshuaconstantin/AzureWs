package router

import (
	"github.com/gorilla/mux"

	"go-postgres-crud/controller"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/buku", controller.AmbilSemuaBuku).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/buku/{id}", controller.AmbilBuku).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/buku", controller.TmbhBuku).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/buku/{id}", controller.UpdateBuku).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/buku/{id}", controller.HapusBuku).Methods("DELETE", "OPTIONS")

	// User_Login API
	router.HandleFunc("/api/users", controller.GetAllUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{id}", controller.GetSnglUsr).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/add_user", controller.InsrtNewUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", controller.UpdtUserPsswd).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/user/{id}", controller.DltUsr).Methods("DELETE", "OPTIONS")

	return router
}
