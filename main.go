package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var rpcURL = "http://localhost:16534"
var explorerURL = "http://localhost:3000"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allow all methods
		w.Header().Set("Access-Control-Allow-Headers", "*")                               // Allow all headers

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleRPCProxy(w http.ResponseWriter, r *http.Request) {
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

func handleExplorerProxy(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(explorerURL, "application/json", bytes.NewBuffer(body))
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
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", handleRPCProxy)
	mux.HandleFunc("/explorer", handleExplorerProxy)

	handler := corsMiddleware(mux)

	fmt.Println("RPC Proxy running on port 1992...")
	http.ListenAndServe(":1992", handler)
}
