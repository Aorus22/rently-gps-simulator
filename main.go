package main

import (
	"log"
	"net/http"

	c "rently-gps/controllers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
	godotenv.Load()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}).Methods("GET")
	r.HandleFunc("/track/{vehicleid}", c.HandleVehicleTracking)
	r.HandleFunc("/start-track/{vehicleid}", c.HttpStartTracker).Methods("GET")
	r.HandleFunc("/stop-track/{vehicleid}", c.HttpStopTracker).Methods("GET")
	r.HandleFunc("/list-track", c.HttpListTrackers).Methods("GET")
	r.HandleFunc("/list-track/{vehicleid}", c.HttpCheckTrackerStatus).Methods("GET")
	r.HandleFunc("/send-email", c.SendEmail).Methods("POST")

	log.Println("Server jalan di :4040")
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(r)
	log.Fatal(http.ListenAndServe(":4040", handler))
}
