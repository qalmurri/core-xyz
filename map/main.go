package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const mapDataFile = "map_nodes.json"
var mapNodes [][]int

func loadMapData() {
	fileData, err := os.ReadFile(mapDataFile)
	if err != nil {
		log.Fatalf("[ERROR] File %s tidak ditemukan: %v\n", mapDataFile, err)
	}

	// Parsing JSON langsung ke format array 2 Dimensi Integer
	err = json.Unmarshal(fileData, &mapNodes)
	if err != nil {
		log.Fatalf("[ERROR] Format JSON tidak valid: %v\n", err)
	}
	fmt.Printf("[MAP API] Berhasil memuat %d koordinat node ke memori!\n", len(mapNodes))
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getNodesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	// SEMENTARA: Hilangkan filter bounding box
	// Langsung kirim seluruh 1.000 data mapNodes ke frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapNodes)
}

func main() {
	loadMapData()
	http.HandleFunc("/api/v1/map/nodes", getNodesHandler)
	
	// Dijalankan di port terpisah agar tidak bentrok dengan Backend 2
	fmt.Println("Backend 1 (Map Render API) berjalan di port 8081")
	http.ListenAndServe(":8081", nil)
}
