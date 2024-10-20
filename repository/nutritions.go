package repository

import (
	"context"
	sqlc "shorvath/nutrition-tracker/generated"
)

func ConvertNutrition(nutritionSqlc sqlc.Nutrition_sqlc) *Nutrition {
	return &Nutrition{
		ID:                   nutritionSqlc.ID,
		CaloriesPer100g:      nutritionSqlc.CaloriesPer100g,
		FatsPer100g:          nutritionSqlc.FatsPer100g,
		FatsSaturatedPer100g: nutritionSqlc.FatsSaturatedPer100g,
		CarbsPer100g:         nutritionSqlc.CarbsPer100g,
		CarbsSugarPer100g:    nutritionSqlc.CarbsSugarPer100g,
		ProteinsPer100g:      nutritionSqlc.ProteinsPer100g,
		SaltPer100g:          nutritionSqlc.SaltPer100g,
	}
}

func CreateNutrition(ctx context.Context, nutrition Nutrition) (*Nutrition, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(db)
	nutritionSqlc, err := queries.CreateNutrition(ctx, sqlc.CreateNutritionParams{
		CaloriesPer100g:      nutrition.CaloriesPer100g,
		FatsPer100g:          nutrition.FatsPer100g,
		FatsSaturatedPer100g: nutrition.FatsSaturatedPer100g,
		CarbsPer100g:         nutrition.CarbsPer100g,
		CarbsSugarPer100g:    nutrition.CarbsSugarPer100g,
		ProteinsPer100g:      nutrition.ProteinsPer100g,
		SaltPer100g:          nutrition.SaltPer100g,
	})
	if err != nil {
		return nil, err
	}
	return ConvertNutrition(nutritionSqlc), nil
}
