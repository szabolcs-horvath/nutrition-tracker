package repository

import (
	"context"
	"database/sql"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/szabolcs-horvath/nutrition-tracker/generated"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
)

type Item struct {
	ID                     int64
	Name                   string
	Owner                  *User
	Language               *Language
	Liquid                 bool
	DefaultPortion         *Portion
	CaloriesPer100         float64
	FatsPer100             float64
	FatsSaturatedPer100    *float64
	CarbsPer100            float64
	CarbsSugarPer100       *float64
	CarbsSlowReleasePer100 *float64
	CarbsFastReleasePer100 *float64
	ProteinsPer100         float64
	SaltPer100             *float64
	Portions               []*Portion
}

func convertItem(itemSqlc *sqlc.Item_sqlc) *Item {
	return &Item{
		ID:                     itemSqlc.ID,
		Name:                   itemSqlc.Name,
		Liquid:                 itemSqlc.Liquid,
		CaloriesPer100:         itemSqlc.CaloriesPer100,
		FatsPer100:             itemSqlc.FatsPer100,
		FatsSaturatedPer100:    itemSqlc.FatsSaturatedPer100,
		CarbsPer100:            itemSqlc.CarbsPer100,
		CarbsSugarPer100:       itemSqlc.CarbsSugarPer100,
		CarbsSlowReleasePer100: itemSqlc.CarbsSlowReleasePer100,
		CarbsFastReleasePer100: itemSqlc.CarbsFastReleasePer100,
		ProteinsPer100:         itemSqlc.ProteinsPer100,
		SaltPer100:             itemSqlc.SaltPer100,
	}
}

func ListItems(ctx context.Context) ([]*Item, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*Item, len(list))
	for i, n := range list {
		result[i] = convertItem(&n.ItemSqlc)
		if &n.ItemsUsersView != nil {
			result[i].Owner = convertUser(ItemsUsersViewWrapper{n.ItemsUsersView})
		}
		result[i].Language = convertLanguage(LanguageSqlcWrapper{n.LanguageSqlc})
		result[i].DefaultPortion = convertPortion(&n.PortionSqlc)
		if result[i].Owner != nil {
			result[i].Portions, err = ListPortionsForItemAndUser(ctx, result[i].ID, &result[i].Owner.ID)
		} else {
			result[i].Portions, err = ListPortionsForItemAndUser(ctx, result[i].ID, nil)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func FindItemById(ctx context.Context, id int64) (*Item, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	row, err := queries.FindItemById(ctx, id)
	if err != nil {
		return nil, err
	}
	item := convertItem(&row.ItemSqlc)
	item.Owner = convertUser(ItemsUsersViewWrapper{row.ItemsUsersView})
	item.Language = convertLanguage(LanguageSqlcWrapper{row.LanguageSqlc})
	item.DefaultPortion = convertPortion(&row.PortionSqlc)
	if item.Owner != nil {
		item.Portions, err = ListPortionsForItemAndUser(ctx, item.ID, &item.Owner.ID)
	} else {
		item.Portions, err = ListPortionsForItemAndUser(ctx, item.ID, nil)
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

func SearchItemsByNameAndUser(ctx context.Context, userId int64, query string) ([]*Item, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.SearchItemsByNameAndUser(ctx, sqlc.SearchItemsByNameAndUserParams{
		OwnerID: &userId,
		Query:   &query,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*Item, len(list))
	for i, n := range list {
		result[i] = convertItem(&n.ItemSqlc)
		if &n.ItemsUsersView != nil {
			result[i].Owner = convertUser(ItemsUsersViewWrapper{n.ItemsUsersView})
		}
		result[i].Language = convertLanguage(LanguageSqlcWrapper{n.LanguageSqlc})
		result[i].DefaultPortion = convertPortion(&n.PortionSqlc)
		if result[i].Owner != nil {
			result[i].Portions, err = ListPortionsForItemAndUser(ctx, result[i].ID, &result[i].Owner.ID)
		} else {
			result[i].Portions, err = ListPortionsForItemAndUser(ctx, result[i].ID, nil)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type CreateItemRequest struct {
	Name                   string   `json:"name"`
	OwnerID                *int64   `json:"owner_id"`
	LanguageID             int64    `json:"language_id"`
	Liquid                 bool     `json:"liquid"`
	DefaultPortionID       int64    `json:"default_portion_id"`
	CaloriesPer100         float64  `json:"calories_per_100"`
	FatsPer100             float64  `json:"fats_per_100"`
	FatsSaturatedPer100    *float64 `json:"fats_saturated_per_100"`
	CarbsPer100            float64  `json:"carbs_per_100"`
	CarbsSugarPer100       *float64 `json:"carbs_sugar_per_100"`
	CarbsSlowReleasePer100 *float64 `json:"carbs_slow_release_per_100"`
	CarbsFastReleasePer100 *float64 `json:"carbs_fast_release_per_100"`
	ProteinsPer100         float64  `json:"proteins_per_100"`
	SaltPer100             *float64 `json:"salt_per_100"`
	PortionIDs             []int64  `json:"portion_ids"`
}

func CreateItem(ctx context.Context, item CreateItemRequest) (*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return createItem(ctx, db, item)
}

func createItem(ctx context.Context, db *sql.DB, item CreateItemRequest) (*Item, error) {
	var result sqlc.Item_sqlc
	err := DoInTransaction(ctx, db, func(childCtx context.Context, queries *sqlc.Queries) error {
		itemSqlc, err := queries.CreateItem(childCtx, sqlc.CreateItemParams{
			Name:                   item.Name,
			OwnerID:                item.OwnerID,
			LanguageID:             item.LanguageID,
			Liquid:                 item.Liquid,
			DefaultPortionID:       item.DefaultPortionID,
			CaloriesPer100:         item.CaloriesPer100,
			FatsPer100:             item.FatsPer100,
			FatsSaturatedPer100:    item.FatsSaturatedPer100,
			CarbsPer100:            item.CarbsPer100,
			CarbsSugarPer100:       item.CarbsSugarPer100,
			CarbsSlowReleasePer100: item.CarbsSlowReleasePer100,
			CarbsFastReleasePer100: item.CarbsFastReleasePer100,
			ProteinsPer100:         item.ProteinsPer100,
			SaltPer100:             item.SaltPer100,
		})
		if err != nil {
			return err
		}

		portionIDsSet := mapset.NewSet[int64](item.PortionIDs...)
		portionIDsSet.Add(item.DefaultPortionID)
		err = addRecordsToItemsPortionsJoiningTable(childCtx, db, itemSqlc.ID, mapset.NewSet[int64](), portionIDsSet)
		if err != nil {
			return err
		}

		result = itemSqlc
		return nil
	})
	if err != nil {
		return nil, err
	}
	return convertItem(&result), nil
}

func addRecordsToItemsPortionsJoiningTable(ctx context.Context, db *sql.DB, itemId int64, existingPortionIds, newPortionIDs mapset.Set[int64]) error {
	err := DoInTransaction(ctx, db, func(childCtx context.Context, queries *sqlc.Queries) error {
		for _, portionId := range newPortionIDs.ToSlice() {
			if !existingPortionIds.Contains(portionId) {
				isDefault, err := queries.IsDefaultPortion(childCtx, portionId)
				if err != nil {
					return err
				}
				if !isDefault {
					if err = queries.CreateItemsPortionsJoiningTableRecord(childCtx, sqlc.CreateItemsPortionsJoiningTableRecordParams{
						ItemID:    itemId,
						PortionID: portionId,
					}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateMultipleItems(ctx context.Context, items []CreateItemRequest) ([]*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var result = make([]*Item, len(items))
	err = DoInTransaction(ctx, db, func(childCtx context.Context, _ *sqlc.Queries) error {
		for i, item := range items {
			createdItem, createErr := createItem(childCtx, db, item)
			if createErr != nil {
				return createErr
			}
			result[i] = createdItem
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

type UpdateItemRequest struct {
	ID                     int64    `json:"id"`
	Name                   string   `json:"name"`
	OwnerID                *int64   `json:"owner_id"`
	LanguageID             int64    `json:"language_id"`
	Liquid                 bool     `json:"liquid"`
	DefaultPortionID       int64    `json:"default_portion_id"`
	CaloriesPer100         float64  `json:"calories_per_100"`
	FatsPer100             float64  `json:"fats_per_100"`
	FatsSaturatedPer100    *float64 `json:"fats_saturated_per_100"`
	CarbsPer100            float64  `json:"carbs_per_100"`
	CarbsSugarPer100       *float64 `json:"carbs_sugar_per_100"`
	CarbsSlowReleasePer100 *float64 `json:"carbs_slow_release_per_100"`
	CarbsFastReleasePer100 *float64 `json:"carbs_fast_release_per_100"`
	ProteinsPer100         float64  `json:"proteins_per_100"`
	SaltPer100             *float64 `json:"salt_per_100"`
	PortionIDs             []int64  `json:"portion_ids"`
}

func UpdateItem(ctx context.Context, item UpdateItemRequest) (*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var result sqlc.Item_sqlc
	err = DoInTransaction(ctx, db, func(childCtx context.Context, queries *sqlc.Queries) error {
		itemSqlc, UpdateErr := queries.UpdateItem(ctx, sqlc.UpdateItemParams{
			Name:                   item.Name,
			OwnerID:                item.OwnerID,
			LanguageID:             item.LanguageID,
			Liquid:                 item.Liquid,
			DefaultPortionID:       item.DefaultPortionID,
			CaloriesPer100:         item.CaloriesPer100,
			FatsPer100:             item.FatsPer100,
			FatsSaturatedPer100:    item.FatsSaturatedPer100,
			CarbsPer100:            item.CarbsPer100,
			CarbsSugarPer100:       item.CarbsSugarPer100,
			CarbsSlowReleasePer100: item.CarbsSlowReleasePer100,
			CarbsFastReleasePer100: item.CarbsFastReleasePer100,
			ProteinsPer100:         item.ProteinsPer100,
			SaltPer100:             item.SaltPer100,
			ID:                     item.ID,
		})
		if UpdateErr != nil {
			return UpdateErr
		}

		portionsOfItem, portionsErr := ListNonDefaultPortionsForItem(childCtx, item.ID)
		if portionsErr != nil {
			return portionsErr
		}
		existingPortionIDs := mapset.NewSet(util.Map(portionsOfItem, func(p *Portion) int64 { return p.ID })...)
		newPortionIDs := mapset.NewSet[int64](item.PortionIDs...)
		newPortionIDs.Add(item.DefaultPortionID)
		err = addRecordsToItemsPortionsJoiningTable(childCtx, db, itemSqlc.ID, existingPortionIDs, newPortionIDs)
		result = itemSqlc
		return nil
	})

	return convertItem(&result), nil
}

func DeleteItem(ctx context.Context, id int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.DeleteItem(ctx, id); err != nil {
		return err
	}
	return nil
}

func GetOwnerIdByItemId(ctx context.Context, itemId int64) (*int64, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	ownerId, err := queries.GetOwnerIdByItemId(ctx, itemId)
	if err != nil {
		return nil, err
	}
	return ownerId, nil
}
