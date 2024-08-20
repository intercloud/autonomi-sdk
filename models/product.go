package models

import "time"

type ProviderType string

const (
	ProviderTypeIntercloud ProviderType = "InterCloud"
	ProviderTypeMegaport   ProviderType = "MEGAPORT"
	ProviderTypeEquinix    ProviderType = "EQUINIX"
)

func (pt ProviderType) String() string {
	return string(pt)
}

type Product struct {
	Provider  ProviderType `json:"provider"`
	Duration  int          `json:"duration"`
	Location  string       `json:"location"`
	Bandwidth int          `json:"bandwidth"`
	Date      time.Time    `json:"date"`
	PriceNRC  int          `json:"priceNrc"`
	PriceMRC  int          `json:"priceMrc"`
	CostNRC   int          `json:"costNrc"`
	CostMRC   int          `json:"costMrc"`
	SKU       string       `json:"sku"`
}
