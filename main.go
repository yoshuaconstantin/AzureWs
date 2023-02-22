package main

import (
	"AzureWS/router"
	"AzureWS/websocketstruct"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Server dijalankan pada port 8080...")

	//Quit infinite loop if os give signal to quit
	quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
	//End

	// Auto Delete Community chat after 30 days old chat
	go websocketstruct.DeleteOldMessage()

	log.Fatal(http.ListenAndServe(":8080", r))
}
