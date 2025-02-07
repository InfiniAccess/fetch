package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"fetch/internal/api"
	"fetch/internal/storage"
)

func main() {
	store := storage.NewStore()
	handler := api.NewHandler(store)
	mux := http.NewServeMux()

	mux.HandleFunc("/receipts/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fmt.Printf("Received request: %s %s\n", r.Method, path)

		switch {
		case path == "/receipts/process":
			handler.ProcessReceipt(w, r)
		case strings.HasSuffix(path, "/points"):
			handler.GetPoints(w, r)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
