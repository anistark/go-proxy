package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var rpcURL = "http://localhost:16534"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allow all methods
		w.Header().Set("Access-Control-Allow-Headers", "*")                               // Allow all headers

		// If it's a preflight request, return immediately
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Forward request to the actual RPC
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to connect to RPC", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response back to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/proxy", handleProxy) // Single endpoint for the proxy

	// Wrap with CORS middleware
	handler := corsMiddleware(mux)

	fmt.Println("RPC Proxy running on port 1992...")
	http.ListenAndServe(":1992", handler)
}
