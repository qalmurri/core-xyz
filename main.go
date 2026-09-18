package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const dataFileName = "nodes.json" // Nama file JSON penyimpan data node

// Patrol menyimpan data kecepatan, radius, warna, dan kemiringan orbit armada
type Patrol struct {
	Speed  float64 `json:"speed"`
	Radius float64 `json:"radius"`
	Color  string  `json:"color"`
	Tilt   float64 `json:"tilt"`
}

// Node menyimpan entitas planet/base/tambang beserta atribut lokasi XYZ dan data detailnya
type Node struct {
	ID      string   `json:"id"`
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type"`
	Level   int      `json:"level"`
	X       int      `json:"x"`
	Y       int      `json:"y"`
	Z       int      `json:"z"`
	Patrols []Patrol `json:"patrols"`
	Owner   string   `json:"owner,omitempty"`
	Points  int      `json:"points,omitempty"`
}

var (
	nodesData []Node    // Memory storage untuk menampung seluruh data node game
	once      sync.Once // Memastikan pemanggilan pembacaan file hanya dieksekusi 1 kali saat server start
)

// loadNodesData khusus membaca data dari file nodes.json dan akan menghentikan server jika file tidak ditemukan
func loadNodesData() {
	fmt.Printf("[STORAGE] Membaca data node dari file %s...\n", dataFileName)
	fileData, err := os.ReadFile(dataFileName) // Membaca isi file nodes.json dari disk
	if err != nil {
		log.Fatalf("[ERROR] File %s tidak ditemukan atau gagal dibaca! Pastikan file tersebut ada di direktori yang sama: %v\n", dataFileName, err)
	}

	err = json.Unmarshal(fileData, &nodesData) // Parsing isi file JSON ke dalam slice nodesData
	if err != nil {
		log.Fatalf("[ERROR] Format struktur JSON pada file %s tidak valid: %v\n", dataFileName, err)
	}

	fmt.Printf("[STORAGE] Berhasil memuat %d node dari %s!\n", len(nodesData), dataFileName)
}

// enableCORS mengizinkan request dari frontend beda port/domain
func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// getNodesHandler memfilter dan mengembalikan node ringan di dalam bounding box area pandang
func getNodesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return // Return langsung jika browser mengirim preflight request
	}

	// Membaca parameter batas area XYZ dari query string URL
	minX, _ := strconv.Atoi(r.URL.Query().Get("minX"))
	maxX, _ := strconv.Atoi(r.URL.Query().Get("maxX"))
	minY, _ := strconv.Atoi(r.URL.Query().Get("minY"))
	maxY, _ := strconv.Atoi(r.URL.Query().Get("maxY"))
	minZ, _ := strconv.Atoi(r.URL.Query().Get("minZ"))
	maxZ, _ := strconv.Atoi(r.URL.Query().Get("maxZ"))

	// Set rentang default (-10 sampai 10) jika parameter query kosong
	if maxX == 0 && minX == 0 {
		minX, maxX, minY, maxY, minZ, maxZ = -10, 10, -10, 10, -10, 10
	}

	var filtered []Node
	for _, node := range nodesData {
		// Filter node yang posisinya masuk dalam batas bounding box kamera
		if node.X >= minX && node.X <= maxX &&
			node.Y >= minY && node.Y <= maxY &&
			node.Z >= minZ && node.Z <= maxZ {
			filtered = append(filtered, Node{
				ID:      node.ID,
				Type:    node.Type,
				Level:   node.Level,
				X:       node.X,
				Y:       node.Y,
				Z:       node.Z,
				Patrols: node.Patrols,
			}) // Mengirimkan payload ringan tanpa field Name, Owner, dan Points
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered) // Serialize slice node ke format JSON
}

// getNodeDetailHandler mengembalikan detail lengkap node saat planet tertentu diklik
func getNodeDetailHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	id := r.URL.Query().Get("id") // Mengambil ID node dari query string
	for _, node := range nodesData {
		if node.ID == id {
			// Mencetak console log di terminal saat ada player yang mengklik node ini
			fmt.Printf("[CLICK LOG] [%s] Player klik Node ID: %s | Nama: %s | Koordinat: (%d, %d, %d) | Pemilik: %s | Poin: %d\n",
				time.Now().Format("15:04:05"), node.ID, node.Name, node.X, node.Y, node.Z, node.Owner, node.Points)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(node) // Mengembalikan data node lengkap jika ditemukan
			return
		}
	}

	// Mencetak console log peringatan jika ID node yang diklik tidak ditemukan
	fmt.Printf("[CLICK LOG] [%s] Player mencoba klik Node ID: %s tetapi tidak ditemukan!\n", time.Now().Format("15:04:05"), id)
	http.Error(w, "Objek tidak ditemukan", http.StatusNotFound)
}

func main() {
	once.Do(loadNodesData) // Membaca file nodes.json saat server pertama kali dinyalakan

	http.HandleFunc("/api/v1/map/nodes", getNodesHandler)     // Endpoint untuk list node di area
	http.HandleFunc("/api/v1/map/node", getNodeDetailHandler) // Endpoint untuk detail 1 node

	fmt.Println("Server Go berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil) // Menjalankan HTTP server pada port 8080
}
