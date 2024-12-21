package repository

import (
	"context"
	"database/sql"
	"github.com/szabolcs-horvath/nutrition-tracker/generated"
)

type Item struct {
	ID                     int64
	Name                   string
	Icon                   []byte
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
		Icon:                   itemSqlc.Icon,
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
		if &n.UserSqlc != nil {
			result[i].Owner = convertUser(UserSqlcWrapper{n.UserSqlc})
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

func CreateItem(ctx context.Context, item *Item) (*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return createItem(ctx, db, item)
}

func createItem(ctx context.Context, db *sql.DB, item *Item) (*Item, error) {
	queries := sqlc.New(db)
	itemSqlc, err := queries.CreateItem(ctx, sqlc.CreateItemParams{
		Name:                   item.Name,
		Icon:                   item.Icon,
		OwnerID:                &item.Owner.ID,
		LanguageID:             item.Language.ID,
		Liquid:                 item.Liquid,
		DefaultPortionID:       item.DefaultPortion.ID,
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
		return nil, err
	}
	return convertItem(&itemSqlc), nil
}

func CreateMultipleItems(ctx context.Context, items []*Item) ([]*Item, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var result = make([]*Item, len(items))
	err = DoInTransaction(ctx, db, func(*sqlc.Queries) error {
		for i, item := range items {
			createdItem, createErr := createItem(ctx, db, item)
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

func UpdateItem(ctx context.Context, item *Item) (*Item, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	itemSqlc, err := queries.UpdateItem(ctx, sqlc.UpdateItemParams{
		Name:                   item.Name,
		Icon:                   item.Icon,
		OwnerID:                &item.Owner.ID,
		LanguageID:             item.Language.ID,
		Liquid:                 item.Liquid,
		DefaultPortionID:       item.DefaultPortion.ID,
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
	if err != nil {
		return nil, err
	}
	return convertItem(&itemSqlc), nil
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
