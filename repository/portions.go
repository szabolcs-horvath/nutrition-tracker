package repository

import (
	"context"
	"github.com/szabolcs-horvath/nutrition-tracker/generated"
)

type Portion struct {
	ID            int64
	Name          string
	Owner         *User
	Language      *Language
	Liquid        bool
	WeightInGrams *float64
	VolumeInMls   *float64
}

func (p Portion) getMultiplier() float64 {
	if p.Liquid {
		return *p.VolumeInMls
	} else {
		return *p.WeightInGrams
	}
}

func convertPortion(portion *sqlc.Portion_sqlc) *Portion {
	return &Portion{
		ID:            portion.ID,
		Name:          portion.Name,
		Liquid:        portion.Liquid,
		WeightInGrams: portion.WeigthInGrams,
		VolumeInMls:   portion.VolumeInMl,
	}
}

func ListPortionsForItemAndUser(ctx context.Context, itemId int64, ownerId *int64) ([]*Portion, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.FindPortionsForItemAndUser(ctx, sqlc.FindPortionsForItemAndUserParams{
		ItemID:  itemId,
		OwnerID: ownerId,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*Portion, len(list))
	for i, p := range list {
		result[i] = convertPortion(&p.PortionSqlc)
	}
	return result, nil
}

func FindPortionById(ctx context.Context, id int64) (*Portion, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	row, err := queries.FindPortionById(ctx, id)
	if err != nil {
		return nil, err
	}
	portion := convertPortion(&row.PortionSqlc)
	portion.Owner = convertUser(PortionsUsersViewWrapper{row.PortionsUsersView})
	portion.Language = convertLanguage(PortionsLanguagesViewWrapper{row.PortionsLanguagesView})
	return portion, nil
}

func CreatePortion(ctx context.Context, portion *Portion) (*Portion, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	portionSqlc, err := queries.CreatePortion(ctx, sqlc.CreatePortionParams{
		Name:          portion.Name,
		OwnerID:       &portion.Owner.ID,
		LanguageID:    &portion.Language.ID,
		Liquid:        portion.Liquid,
		WeigthInGrams: portion.WeightInGrams,
		VolumeInMl:    portion.VolumeInMls,
	})
	if err != nil {
		return nil, err
	}
	return convertPortion(&portionSqlc), nil
}

func UpdatePortion(ctx context.Context, portion *Portion) (*Portion, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	portionSqlc, err := queries.UpdatePortion(ctx, sqlc.UpdatePortionParams{
		Name:          portion.Name,
		OwnerID:       &portion.Owner.ID,
		LanguageID:    &portion.Language.ID,
		Liquid:        portion.Liquid,
		WeigthInGrams: portion.WeightInGrams,
		VolumeInMl:    portion.VolumeInMls,
		ID:            portion.ID,
	})
	if err != nil {
		return nil, err
	}
	return convertPortion(&portionSqlc), nil
}

func DeletePortion(ctx context.Context, id int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.DeletePortion(ctx, id); err != nil {
		return err
	}
	return nil
}
