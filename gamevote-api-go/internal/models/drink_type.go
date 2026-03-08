package models

// DrinkType represents a preset drink like "Beer 500ml" or "Vodka Shot"
type DrinkType struct {
	ID             string  `json:"id,omitempty" surreal:"id,omitempty"`
	Name           string  `json:"name" surreal:"name"`
	VolumeMl       int     `json:"volumeMl" surreal:"volumeMl"`
	AlcoholPercent float64 `json:"alcoholPercent" surreal:"alcoholPercent"` // Percentage, e.g., 5.0 for 5%
	BeerEquivalent float64 `json:"beerEquivalent" surreal:"beerEquivalent"` // How many beers this equates to
	UnitName       string  `json:"unitName" surreal:"unitName"`             // e.g. "Halber Liter", "Shot", "Achtel"
}
