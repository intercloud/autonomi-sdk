package models

import "time"

type Product struct {
	Provider  string    `json:"provider"`
	Duration  int       `json:"duration"`
	Location  string    `json:"location"`
	Bandwidth int       `json:"bandwidth"`
	Date      time.Time `json:"date"`
	PriceNRC  int       `json:"priceNrc"`
	PriceMRC  int       `json:"priceMrc"`
	CostNRC   int       `json:"costNrc"`
	CostMRC   int       `json:"costMrc"`
	SKU       string    `json:"sku"`
}
