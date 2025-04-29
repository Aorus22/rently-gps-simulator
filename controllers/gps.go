package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	u "rently-gps/utils"
	"github.com/gorilla/mux"
)

func HandleVehicleTracking(res http.ResponseWriter, req *http.Request) {
	vehicleID := mux.Vars(req)["vehicleid"]

	u.Mu.Lock()
	_, active := u.ActiveTracks[vehicleID]
	u.Mu.Unlock()

	if !active {
		log.Printf("[%s] Tracker non-aktif, mengirim posisi terakhir", vehicleID)

		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Cache-Control", "no-cache")
		res.Header().Set("Connection", "keep-alive")

		pos := u.GetPosition(vehicleID)

		if pos.Latitude == 0 || pos.Longitude == 0 || pos.Timestamp == 0 {
			http.Error(res, "Posisi terakhir tidak tersedia", http.StatusNotFound)
			return
		}

		data, _ := json.Marshal(pos)
		res.WriteHeader(http.StatusPartialContent)
		res.Write(data)
		return
	}

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	log.Printf("[%s] SSE terhubung", vehicleID)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pos := u.GetPosition(vehicleID)
		data, _ := json.Marshal(pos)
		_, err := res.Write([]byte("data: " + string(data) + "\n\n"))
		if err != nil {
			log.Printf("[%s] SSE terputus", vehicleID)
			return
		}
		res.(http.Flusher).Flush()
	}
}

func HttpStartTracker(res http.ResponseWriter, req *http.Request) {
	vehicleID := mux.Vars(req)["vehicleid"]
	u.StartTracker(vehicleID)

	resp := map[string]string{
		"message":    "Tracker dimulai",
		"vehicle_id": vehicleID,
		"status":     "active",
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
}

func HttpStopTracker(res http.ResponseWriter, req *http.Request) {
	vehicleID := mux.Vars(req)["vehicleid"]
	u.StopTracker(vehicleID)

	resp := map[string]string{
		"message":    "Tracker dihentikan",
		"vehicle_id": vehicleID,
		"status":     "nonactive",
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
}

func HttpListTrackers(res http.ResponseWriter, req *http.Request) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	var list []string
	for vehicleID := range u.ActiveTracks {
		list = append(list, vehicleID)
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(list)
}

func HttpCheckTrackerStatus(res http.ResponseWriter, req *http.Request) {
	vehicleID := mux.Vars(req)["vehicleid"]

	u.Mu.Lock()
	_, active := u.ActiveTracks[vehicleID]
	u.Mu.Unlock()

	status := "nonactive"
	if active {
		status = "active"
	}

	resp := map[string]string{
		"vehicle_id": vehicleID,
		"status":     status,
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(resp)
}
