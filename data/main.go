package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const detailDataFile = "node_details.json"

// NodeDetail hanya menyimpan struktur teks dan informasi UI
type NodeDetail struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Points      int    `json:"points"`
	TypeName    string `json:"type_name"`
	Description string `json:"description"`
}

var nodeDetails []NodeDetail

func loadDetailData() {
	fileData, err := os.ReadFile(detailDataFile)
	if err != nil {
		log.Fatalf("[ERROR] File %s tidak ditemukan: %v\n", detailDataFile, err)
	}

	err = json.Unmarshal(fileData, &nodeDetails)
	if err != nil {
		log.Fatalf("[ERROR] Format JSON tidak valid: %v\n", err)
	}
	fmt.Printf("[DETAIL API] Berhasil memuat %d profil node ke memori!\n", len(nodeDetails))
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// getNodeDetailHandler memproses request berdasarkan query parameter "id"
func getNodeDetailHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	idStr := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Format ID tidak valid", http.StatusBadRequest)
		return
	}

	// Mencari data detail yang cocok dengan ID
	for _, node := range nodeDetails {
		if node.ID == idInt {
			fmt.Printf("[CLICK LOG] [%s] Player memuat detail Node ID: %d | Nama: %s\n",
				time.Now().Format("15:04:05"), node.ID, node.Name)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(node)
			return
		}
	}

	http.Error(w, "Objek tidak ditemukan", http.StatusNotFound)
}

func main() {
	loadDetailData()

	// Endpoint khusus untuk detail node
	http.HandleFunc("/api/v1/map/node", getNodeDetailHandler)

	fmt.Println("Backend 2 (Detail API) berjalan di port 8080")
	http.ListenAndServe(":8080", nil)
}
