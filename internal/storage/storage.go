package storage

import (
	"fmt"
	"sync"

	"fetch/internal/functions"
	"fetch/internal/models"

	"github.com/google/uuid"
)

// Store manages thread-safe storage of receipts and their points
type Store struct {
	sync.RWMutex
	receipts   map[string]*models.StoredReceipt
	calculator *functions.ReceiptCalculator
}

// NewStore creates a new Store instance with initialized storage
func NewStore() *Store {
	return &Store{
		receipts:   make(map[string]*models.StoredReceipt),
		calculator: functions.NewReceiptCalculator(),
	}
}

// CalculateAndSaveReceipt calculates points and stores the receipt
func (s *Store) CalculateAndSaveReceipt(receipt *models.Receipt) (string, int) {
	points := s.calculator.CalculatePoints(receipt)
	id := s.SaveReceipt(receipt, points)
	return id, points
}

// SaveReceipt stores a receipt and returns its ID
func (s *Store) SaveReceipt(receipt *models.Receipt, points int) string {
	s.Lock()
	defer s.Unlock()

	id := uuid.New().String()
	s.receipts[id] = &models.StoredReceipt{
		Receipt: *receipt,
		Points:  points,
	}
	return id
}

// GetReceipt retrieves a receipt by ID
func (s *Store) GetReceipt(id string) (*models.StoredReceipt, error) {
	s.RLock()
	defer s.RUnlock()

	if receipt, ok := s.receipts[id]; ok {
		return receipt, nil
	}
	return nil, fmt.Errorf("receipt not found")
}

// GetAllReceipts returns all stored receipts
func (s *Store) GetAllReceipts() map[string]*models.StoredReceipt {
	s.RLock()
	defer s.RUnlock()

	receipts := make(map[string]*models.StoredReceipt)
	for id, receipt := range s.receipts {
		receipts[id] = receipt
	}
	return receipts
}

// getReceiptIDs returns all receipt IDs
func (s *Store) getReceiptIDs() []string {
	s.RLock()
	defer s.RUnlock()

	ids := make([]string, 0, len(s.receipts))
	for id := range s.receipts {
		ids = append(ids, id)
	}
	return ids
}
