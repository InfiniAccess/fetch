package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"fetch/internal/models"
	"fetch/internal/storage"
)

// Handler manages HTTP request handling with storage access
type Handler struct {
	store *storage.Store
}

// NewHandler creates a new Handler instance with the provided storage
func NewHandler(store *storage.Store) *Handler {
	return &Handler{store: store}
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	// sendErrorResponse sends a JSON error response with the given status code
	sendErrorResponse = func(w http.ResponseWriter, message string, status int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(ErrorResponse{Error: message})
	}

	// validateReceipt performs basic validation on the receipt
	validateReceipt = func(receipt *models.Receipt) error {
		if strings.TrimSpace(receipt.Retailer) == "" {
			return fmt.Errorf("retailer name is required")
		}
		if len(receipt.Items) == 0 {
			return fmt.Errorf("at least one item is required")
		}
		
		var sum float64
		for i, item := range receipt.Items {
			if strings.TrimSpace(item.ShortDescription) == "" {
				return fmt.Errorf("item %d: description is required", i+1)
			}
			if item.ParsedPrice <= 0 {
				return fmt.Errorf("item %d: invalid price", i+1)
			}
			sum += item.ParsedPrice
		}
		
		if diff := receipt.ParsedTotal - sum; diff < -0.01 || diff > 0.01 {
			return fmt.Errorf("total %.2f does not match sum of items %.2f", receipt.ParsedTotal, sum)
		}
		
		return nil
	}
)

// ProcessReceipt handles the POST /receipts/process endpoint. It accepts a JSON receipt 
// in the request body and returns a JSON response with the receipt ID and points awarded
func (h *Handler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		sendErrorResponse(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validateReceipt(&receipt); err != nil {
		fmt.Printf("Receipt validation failed: %v\n", err)
		sendErrorResponse(w, fmt.Sprintf("Invalid receipt: %v", err), http.StatusBadRequest)
		return
	}

	id, points := h.store.CalculateAndSaveReceipt(&receipt)
	fmt.Printf("Created receipt with ID: %s and points: %d\n", id, points)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     id,
		"points": points,
	}); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		sendErrorResponse(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}

// GetPoints handles the GET /receipts/{id}/points endpoint. It extracts the receipt ID 
// from the URL path and returns the points awarded for that receipt
func (h *Handler) GetPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id := strings.TrimSuffix(path, "/points")
	
	if id == "" || id == path {
		sendErrorResponse(w, "Invalid receipt ID format", http.StatusBadRequest)
		return
	}

	fmt.Printf("Path: %s, Extracted ID: %s\n", r.URL.Path, id)

	receipt, err := h.store.GetReceipt(id)
	if err != nil {
		fmt.Printf("Error getting receipt: %v\n", err)
		sendErrorResponse(w, "Receipt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]int{
		"points": receipt.Points,
	}); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		sendErrorResponse(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}
