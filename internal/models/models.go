package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// NewItem returns a new Item instance
func NewItem(shortDescription, price string) Item {
	return Item{
		ShortDescription: shortDescription,
		Price:            price,
	}
}

// NewReceipt returns a new Receipt instance
func NewReceipt(retailer, purchaseDate, purchaseTime string, items []Item, total string) Receipt {
	return Receipt{
		Retailer:     retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        items,
		Total:        total,
	}
}

// NewStoredReceipt returns a new StoredReceipt instance
func NewStoredReceipt(receipt Receipt, points int) StoredReceipt {
	return StoredReceipt{
		Receipt: receipt,
		Points:  points,
	}
}

// Item represents a single item on a receipt with its description and price
type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            string  `json:"price"`
	ParsedPrice      float64 `json:"-"`
}

// Receipt represents a complete shopping receipt with purchase details
type Receipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        string    `json:"total"`
	ParsedDate   time.Time `json:"-"`
	ParsedTotal  float64   `json:"-"`
}

// StoredReceipt represents a processed receipt with calculated points
type StoredReceipt struct {
	Receipt Receipt
	Points  int
}

// UnmarshalJSON parses and validates the receipt's date, time, and monetary values
func (r *Receipt) UnmarshalJSON(data []byte) error {
	type TempReceipt Receipt
	var temp TempReceipt

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	dateStr := fmt.Sprintf("%sT%s:00", temp.PurchaseDate, temp.PurchaseTime)
	parsedDate, err := time.Parse("2006-01-02T15:04:05", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date/time format: %v", err)
	}

	total, err := strconv.ParseFloat(temp.Total, 64)
	if err != nil {
		return fmt.Errorf("invalid total: %v", err)
	}

	for i := range temp.Items {
		price, err := strconv.ParseFloat(temp.Items[i].Price, 64)
		if err != nil {
			return fmt.Errorf("invalid price for item %d: %v", i+1, err)
		}
		temp.Items[i].ParsedPrice = price
	}

	*r = Receipt(temp)
	r.ParsedDate = parsedDate
	r.ParsedTotal = total

	return nil
}
