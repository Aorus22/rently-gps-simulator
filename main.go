package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}).Methods("GET")
	r.HandleFunc("/track/{vehicleid}", handleVehicleTracking)
	r.HandleFunc("/start-track/{vehicleid}", httpStartTracker).Methods("GET")
	r.HandleFunc("/stop-track/{vehicleid}", httpStopTracker).Methods("GET")
	r.HandleFunc("/list-track", httpListTrackers).Methods("GET")
	r.HandleFunc("/list-track/{vehicleid}", httpCheckTrackerStatus).Methods("GET")

	log.Println("Server jalan di :4040")
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(r)
	log.Fatal(http.ListenAndServe(":4040", handler))
}
