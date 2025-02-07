package functions

import (
	"math"
	"strings"
	"time"
	"unicode"

	"fetch/internal/models"
)

// ReceiptCalculator is a calculator for receipt points
type ReceiptCalculator struct{}

// NewReceiptCalculator returns a new ReceiptCalculator
func NewReceiptCalculator() *ReceiptCalculator {
	return &ReceiptCalculator{}
}

var (
	// points for round dollar amounts
	roundDollarPoints = 50
	// points for multiples of 0.25
	multipleOf25Points = 25
	// points for odd-numbered days
	oddDayPoints = 6
	// points for purchases between 2:00 PM and 4:00 PM
	between2And4PMPoints = 10
)

// CalculatePoints calculates points for a receipt based on retailer name, total amount, items, date and time
func (rc *ReceiptCalculator) CalculatePoints(receipt *models.Receipt) int {
	points := 0
	points += rc.calculateRetailerNamePoints(receipt.Retailer)
	points += rc.calculateTotalPoints(receipt.ParsedTotal)
	points += rc.calculateItemPoints(receipt.Items)
	points += rc.calculateDatePoints(receipt.ParsedDate)
	points += rc.calculateTimePoints(receipt.ParsedDate)
	return points
}

// calculateRetailerNamePoints counts alphanumeric characters in retailer name
func (rc *ReceiptCalculator) calculateRetailerNamePoints(retailer string) int {
	points := 0
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			points++
		}
	}
	return points
}

// calculateTotalPoints awards points for round dollar amounts and multiples of 0.25
func (rc *ReceiptCalculator) calculateTotalPoints(total float64) int {
	points := 0

	if total == float64(int(total)) {
		points += roundDollarPoints
	}

	if math.Mod(total*100, 25) == 0 {
		points += multipleOf25Points
	}

	return points
}

// calculateItemPoints awards points based on item count and description length
func (rc *ReceiptCalculator) calculateItemPoints(items []models.Item) int {
	points := 0
	points += (len(items) / 2) * 5

	for _, item := range items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 == 0 {
			points += int(math.Ceil(item.ParsedPrice * 0.2))
		}
	}

	return points
}

// calculateDatePoints awards points for odd-numbered days
func (rc *ReceiptCalculator) calculateDatePoints(date time.Time) int {
	if date.Day()%2 == 1 {
		return oddDayPoints
	}
	return 0
}

// calculateTimePoints awards points for purchases between 2:00 PM and 4:00 PM
func (rc *ReceiptCalculator) calculateTimePoints(date time.Time) int {
	hour := date.Hour()
	if hour >= 14 && hour < 16 {
		return between2And4PMPoints
	}
	return 0
}
