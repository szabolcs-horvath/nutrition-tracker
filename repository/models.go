package repository

type Item struct {
	ID        int64
	Name      string
	Nutrition Nutrition
	Icon      []byte
}

type Nutrition struct {
	ID                   int64
	CaloriesPer100g      float64
	FatsPer100g          float64
	FatsSaturatedPer100g *float64
	CarbsPer100g         float64
	CarbsSugarPer100g    *float64
	ProteinsPer100g      float64
	SaltPer100g          *float64
}
