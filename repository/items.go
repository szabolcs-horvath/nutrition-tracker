package repository

import (
	"context"
	sqlc "shorvath/nutrition-tracker/generated"
)

func ConvertItem(itemSqlc sqlc.Item_sqlc, nutritionSqlc sqlc.Nutrition_sqlc) *Item {
	return &Item{
		ID:        itemSqlc.ID,
		Name:      itemSqlc.Name,
		Nutrition: *ConvertNutrition(nutritionSqlc),
		Icon:      itemSqlc.Icon,
	}
}

func ConvertItemWithoutNutrition(itemSqlc sqlc.Item_sqlc) *Item {
	return &Item{
		ID:   itemSqlc.ID,
		Name: itemSqlc.Name,
		Icon: itemSqlc.Icon,
	}
}

func FindItemByIdWithNutrition(ctx context.Context, id int64) (*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(db)
	row, err := queries.FindItemByIdWithNutrition(ctx, id)
	if err != nil {
		return nil, err
	}
	return ConvertItem(row.ItemSqlc, row.NutritionSqlc), nil
}

func CreateItem(ctx context.Context, item Item) (*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	nutrition, err := CreateNutrition(ctx, item.Nutrition)
	if err != nil {
		return nil, err
	}
	itemSqlc, err := queries.CreateItem(ctx, sqlc.CreateItemParams{
		Name:        item.Name,
		NutritionID: nutrition.ID,
		Icon:        item.Icon,
	})
	if err != nil {
		return nil, err
	}
	createdItem := ConvertItemWithoutNutrition(itemSqlc)
	createdItem.Nutrition = *nutrition
	return createdItem, nil
}
