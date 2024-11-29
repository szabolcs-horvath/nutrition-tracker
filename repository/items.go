package repository

import (
	"context"
	sqlc "shorvath/nutrition-tracker/generated"
)

type Item struct {
	ID        int64
	Name      string
	Nutrition *Nutrition
	Icon      []byte
}

func convertItem(itemSqlc *sqlc.Item_sqlc) *Item {
	return &Item{
		ID:   itemSqlc.ID,
		Name: itemSqlc.Name,
		Icon: itemSqlc.Icon,
	}
}

func ListItems(ctx context.Context) ([]*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	queries := sqlc.New(db)
	list, err := queries.ListItemsWithNutritions(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*Item, len(list))
	for i, n := range list {
		result[i] = convertItem(&n.ItemSqlc)
		result[i].Nutrition = convertNutrition(&n.NutritionSqlc)
	}
	return result, nil
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
	item := convertItem(&row.ItemSqlc)
	item.Nutrition = convertNutrition(&row.NutritionSqlc)
	return item, nil
}

func CreateItem(ctx context.Context, item *Item) (*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var createdItem *Item
	err = DoInTransaction(ctx, db, func() error {
		queries := sqlc.New(db)
		nutrition, createNutErr := CreateNutrition(ctx, item.Nutrition)
		if createNutErr != nil {
			return createNutErr
		}
		itemSqlc, createItemErr := queries.CreateItem(ctx, sqlc.CreateItemParams{
			Name:        item.Name,
			NutritionID: nutrition.ID,
			Icon:        item.Icon,
		})
		if createItemErr != nil {
			return createItemErr
		}
		createdItem = convertItem(&itemSqlc)
		createdItem.Nutrition = nutrition
		return nil
	})
	if err != nil {
		return nil, err
	}
	return createdItem, nil
}

func UpdateItem(ctx context.Context, item *Item) (*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var updatedItem *Item
	err = DoInTransaction(ctx, db, func() error {
		queries := sqlc.New(db)
		itemSqlc, sqlErr := queries.UpdateItem(ctx, sqlc.UpdateItemParams{
			ID:   0,
			Name: item.Name,
			Icon: item.Icon,
		})
		if sqlErr != nil {
			return sqlErr
		}
		updatedItem = convertItem(&itemSqlc)
		if item.Nutrition != nil {
			updatedNutrition, nutUpdateErr := UpdateNutrition(ctx, item.Nutrition)
			if nutUpdateErr != nil {
				return nutUpdateErr
			}
			item.Nutrition = updatedNutrition
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedItem, nil
}

func DeleteItem(ctx context.Context, id int64) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	queries := sqlc.New(db)
	if err = queries.DeleteItem(ctx, id); err != nil {
		return err
	}
	return nil
}
