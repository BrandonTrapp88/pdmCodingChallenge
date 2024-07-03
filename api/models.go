package main

import (
	"time"
)

type Part struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Images      []string          `json:"images"`
	SKU         string            `json:"sku"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Attributes  map[string]string `json:"attributes"`
	FitmentData []string          `json:"fitment_data"`
	Location    string            `json:"location"`
	Shipment    ShipmentInfo      `json:"shipment"`
	Metadata    map[string]string `json:"metadata"`
}

type ShipmentInfo struct {
	Weight    float64 `json:"weight"`
	Size      string  `json:"size"`
	Hazardous bool    `json:"hazardous"`
	Fragile   bool    `json:"fragile"`
}
type PartVersion struct {
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Part      Part      `json:"part"`
}
