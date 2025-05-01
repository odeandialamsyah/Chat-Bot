package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Connect ke MongoDB
	err := ConnectDB()
	if err != nil {
		log.Fatalf("Gagal koneksi ke MongoDB: %v", err)
	}

	// Setup handler untuk webhook
	http.HandleFunc("/webhook", webhookHandler)

	// Jalankan server
	fmt.Println("Server berjalan di port 8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
