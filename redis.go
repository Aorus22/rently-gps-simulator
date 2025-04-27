package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_HOST"),
})

func savePosition(vehicleID string, pos Position) {
	data, _ := json.Marshal(pos)
	rdb.Set(context.Background(), "vehicle:"+vehicleID, data, 0)
}

func getPosition(vehicleID string) Position {
	val, err := rdb.Get(context.Background(), "vehicle:"+vehicleID).Result()
	if err != nil {
		log.Printf("[%s] Gagal ambil posisi: %v", vehicleID, err)
		return Position{}
	}
	var pos Position
	json.Unmarshal([]byte(val), &pos)
	return pos
}
