package repository

import (
	"context"
	sqlc "shorvath/nutrition-tracker/generated"
)

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

func convertNutrition(nutritionSqlc *sqlc.Nutrition_sqlc) *Nutrition {
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

func ListNutritions(ctx context.Context) ([]*Nutrition, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(db)
	list, err := queries.ListNutritions(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*Nutrition, len(list))
	for i, n := range list {
		result[i] = convertNutrition(&n)
	}
	return result, nil
}

func FindNutritionById(ctx context.Context, id int64) (*Nutrition, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(db)
	row, err := queries.FindNutritionById(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertNutrition(&row), nil
}

func CreateNutrition(ctx context.Context, nutrition *Nutrition) (*Nutrition, error) {
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
	return convertNutrition(&nutritionSqlc), nil
}

func UpdateNutrition(ctx context.Context, nutrition *Nutrition) (*Nutrition, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var updatedNutrition *Nutrition
	queries := sqlc.New(db)
	nutritionSqlc, err := queries.UpdateNutrition(ctx, sqlc.UpdateNutritionParams{
		ID:                   nutrition.ID,
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
	updatedNutrition = convertNutrition(&nutritionSqlc)
	return updatedNutrition, nil
}

func DeleteNutrition(ctx context.Context, id int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	queries := sqlc.New(db)
	if err = queries.DeleteNutrition(ctx, id); err != nil {
		return err
	}
	return nil
}
