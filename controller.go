package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func handleVehicleTracking(w http.ResponseWriter, r *http.Request) {
	vehicleID := mux.Vars(r)["vehicleid"]

	mu.Lock()
	_, active := activeTracks[vehicleID]
	mu.Unlock()

	if !active {
		log.Printf("[%s] Tracker non-aktif, mengirim posisi terakhir", vehicleID)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		pos := getPosition(vehicleID)

		if pos.Latitude == 0 || pos.Longitude == 0 || pos.Timestamp == 0 {
			http.Error(w, "Posisi terakhir tidak tersedia", http.StatusNotFound)
			return
		}

		data, _ := json.Marshal(pos)
		w.WriteHeader(http.StatusPartialContent)
		w.Write(data)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Printf("[%s] SSE terhubung", vehicleID)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pos := getPosition(vehicleID)
		data, _ := json.Marshal(pos)
		_, err := w.Write([]byte("data: " + string(data) + "\n\n"))
		if err != nil {
			log.Printf("[%s] SSE terputus", vehicleID)
			return
		}
		w.(http.Flusher).Flush()
	}
}

func httpStartTracker(w http.ResponseWriter, r *http.Request) {
	vehicleID := mux.Vars(r)["vehicleid"]
	startTracker(vehicleID)

	resp := map[string]string{
		"message":    "Tracker dimulai",
		"vehicle_id": vehicleID,
		"status":     "active",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func httpStopTracker(w http.ResponseWriter, r *http.Request) {
	vehicleID := mux.Vars(r)["vehicleid"]
	stopTracker(vehicleID)

	resp := map[string]string{
		"message":    "Tracker dihentikan",
		"vehicle_id": vehicleID,
		"status":     "nonactive",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func httpListTrackers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var list []string
	for vehicleID := range activeTracks {
		list = append(list, vehicleID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func httpCheckTrackerStatus(w http.ResponseWriter, r *http.Request) {
	vehicleID := mux.Vars(r)["vehicleid"]

	mu.Lock()
	_, active := activeTracks[vehicleID]
	mu.Unlock()

	status := "nonactive"
	if active {
		status = "active"
	}

	resp := map[string]string{
		"vehicle_id": vehicleID,
		"status":     status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
