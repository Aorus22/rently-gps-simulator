package utils

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type ActiveTracker struct {
	UpdateTicker *time.Ticker
	Done         chan bool
}

type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type Properties struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Source string `json:"source"`
}

var (
	Mu           sync.Mutex
	ActiveTracks = make(map[string]*ActiveTracker)
)

var indonesiaPolygons [][][]float64

func LoadGeoJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var geo GeoJSON
	if err := json.Unmarshal(data, &geo); err != nil {
		return err
	}

	for _, feature := range geo.Features {
		if feature.Geometry.Type == "MultiPolygon" {
			for _, Multipolygon := range feature.Geometry.Coordinates {
				indonesiaPolygons = append(indonesiaPolygons, Multipolygon...)
			}
		}
	}
	return nil
}

func pointInPolygon(lat, lon float64, polygon [][]float64) bool {
	n := len(polygon)
	inside := false

	j := n - 1
	for i := 0; i < n; i++ {
		xi, yi := polygon[i][0], polygon[i][1]
		xj, yj := polygon[j][0], polygon[j][1]

		intersect := ((yi > lat) != (yj > lat)) &&
			(lon < (xj-xi)*(lat-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
		j = i
	}

	return inside
}

func pointInMultiPolygon(lat, lon float64, polygons [][][]float64) bool {
	for _, polygon := range polygons {
		if pointInPolygon(lat, lon, polygon) {
			return true
		}
	}
	return false
}

func StartTracker(vehicleID string) {
	Mu.Lock()
	defer Mu.Unlock()

	if _, exists := ActiveTracks[vehicleID]; exists {
		return
	}

	tracker := &ActiveTracker{
		Done:         make(chan bool),
		UpdateTicker: time.NewTicker(2 * time.Second),
	}
	ActiveTracks[vehicleID] = tracker

	pos := GetPosition(vehicleID)
	if pos.Latitude == 0 && pos.Longitude == 0 {
		SavePosition(vehicleID, Position{-6.200000, 106.816666, time.Now().Unix()})
	}

	go func() {
		for {
			select {
			case <-tracker.UpdateTicker.C:
				pos := GetPosition(vehicleID)

				for {
					latDelta := (rand.Float64() - 0.5) / 500
					lonDelta := (rand.Float64() - 0.5) / 500

					newLat := pos.Latitude + latDelta
					newLon := pos.Longitude + lonDelta

					if pointInMultiPolygon(newLat, newLon, indonesiaPolygons) {
						pos.Latitude = newLat
						pos.Longitude = newLon
						pos.Timestamp = time.Now().Unix()
						break
					}
				}

				SavePosition(vehicleID, pos)

			case <-tracker.Done:
				log.Printf("[%s] SiMulasi dihentikan", vehicleID)
				return
			}
		}
	}()

	log.Printf("[%s] Tracker diMulai", vehicleID)
}

func StopTracker(vehicleID string) {
	Mu.Lock()
	defer Mu.Unlock()

	tracker, exists := ActiveTracks[vehicleID]
	if !exists {
		return
	}

	tracker.UpdateTicker.Stop()
	tracker.Done <- true
	delete(ActiveTracks, vehicleID)

	log.Printf("[%s] Tracker dihentikan", vehicleID)
}

func init() {
	err := LoadGeoJSON("geo_id.json")
	if err != nil {
		log.Fatal(err)
	}
}
