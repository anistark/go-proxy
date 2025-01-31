package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var rpcURL = "http://localhost:16534"

func handleProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to connect to RPC", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	if envURL := os.Getenv("RPC_URL"); envURL != "" {
		rpcURL = envURL
	}

	r := mux.NewRouter()
	r.HandleFunc("/proxy", handleProxy).Methods("POST", "OPTIONS")

	fmt.Println("Proxy server is running on port 1992...")
	http.ListenAndServe(":1992", r)
}
