package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Fungsi untuk menangani kata kunci dan memberikan balasan otomatis
func handleKeywords(message string) string {
	// Ambil kata kunci dari database
	keywords, err := getKeywordsFromDB()
	if err != nil {
		log.Println("Gagal mengambil kata kunci dari database:", err)
		return "Maaf, ada kesalahan dalam sistem 😅"
	}

	// Cek pesan dan beri balasan sesuai kata kunci yang ada di database
	for _, keyword := range keywords {
		if strings.ToLower(message) == strings.ToLower(keyword.Keyword) {
			return keyword.Reply
		}
	}

	// Jika tidak ada kecocokan, balas pesan default
	return "Maaf, saya tidak paham pesanmu 😅"
}

// Fungsi untuk mengirim balasan menggunakan wa.my.id API
func sendReply(number, message string) {
	// Format payload untuk kirim ke wa.my.id API
	payload := map[string]string{
		"api_key":   os.Getenv("WAPI_KEY"),
		"device_id": os.Getenv("DEVICE_ID"),
		"number":    number,
		"message":   message,
	}

	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post("https://wa.my.id/api/send-message", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Gagal kirim balasan: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Balasan terkirim ke %s", number)
}

// Fungsi untuk menangani webhook
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil Secret Code dari environment
	expectedSecret := os.Getenv("WEBHOOKSECRET")
	secretCode := r.Header.Get("X-Secret-Code")

	// Validasi Secret Code
	if secretCode != expectedSecret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Pesan masuk yang diterima
	var msg IncomingMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Menyimpan pesan ke MongoDB dan mengirim balasan
	log.Printf("Pesan dari %s: %s", msg.Number, msg.Message)

	// Proses untuk memberikan balasan otomatis berdasarkan kata kunci
	reply := handleKeywords(msg.Message)

	// Kirim balasan menggunakan wa.my.id API
	go sendReply(msg.Number, reply)

	// Berikan respon 200 OK
	fmt.Fprintln(w, "Pesan diterima")
}
