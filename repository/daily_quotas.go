package repository

import (
	"context"
	sqlc "github.com/szabolcs-horvath/nutrition-tracker/generated"
	"time"
)

type DailyQuotaFromDB interface {
	getId() *int64
	getOwnerID() *int64
	getArchivedDateTime() *time.Time
	getCalories() *float64
	getFats() *float64
	getFatsSaturated() *float64
	getCarbs() *float64
	getCarbsSugar() *float64
	getCarbsSlowRelease() *float64
	getCarbsFastRelease() *float64
	getProteins() *float64
	getSalt() *float64
}

type DailyQuotaSqlcWrapper struct {
	sqlc.DailyQuota_sqlc
}

func (dailyQuota DailyQuotaSqlcWrapper) getId() *int64 {
	return &dailyQuota.ID
}

func (dailyQuota DailyQuotaSqlcWrapper) getOwnerID() *int64 {
	return &dailyQuota.OwnerID
}

func (dailyQuota DailyQuotaSqlcWrapper) getArchivedDateTime() *time.Time {
	return dailyQuota.ArchivedDateTime
}

func (dailyQuota DailyQuotaSqlcWrapper) getCalories() *float64 {
	return dailyQuota.Calories
}

func (dailyQuota DailyQuotaSqlcWrapper) getFats() *float64 {
	return dailyQuota.Fats
}

func (dailyQuota DailyQuotaSqlcWrapper) getFatsSaturated() *float64 {
	return dailyQuota.FatsSaturated
}

func (dailyQuota DailyQuotaSqlcWrapper) getCarbs() *float64 {
	return dailyQuota.Carbs
}

func (dailyQuota DailyQuotaSqlcWrapper) getCarbsSugar() *float64 {
	return dailyQuota.CarbsSugar
}

func (dailyQuota DailyQuotaSqlcWrapper) getCarbsSlowRelease() *float64 {
	return dailyQuota.CarbsSlowRelease
}

func (dailyQuota DailyQuotaSqlcWrapper) getCarbsFastRelease() *float64 {
	return dailyQuota.CarbsFastRelease
}

func (dailyQuota DailyQuotaSqlcWrapper) getProteins() *float64 {
	return dailyQuota.Proteins
}

func (dailyQuota DailyQuotaSqlcWrapper) getSalt() *float64 {
	return dailyQuota.Salt
}

type UsersDailyQuotasViewWrapper struct {
	sqlc.UsersDailyQuotasView
}

func (dailyQuota UsersDailyQuotasViewWrapper) getId() *int64 {
	return dailyQuota.ID
}

func (dailyQuota UsersDailyQuotasViewWrapper) getOwnerID() *int64 {
	return dailyQuota.OwnerID
}

func (dailyQuota UsersDailyQuotasViewWrapper) getArchivedDateTime() *time.Time {
	return dailyQuota.ArchivedDateTime
}

func (dailyQuota UsersDailyQuotasViewWrapper) getCalories() *float64 {
	return dailyQuota.Calories
}

func (dailyQuota UsersDailyQuotasViewWrapper) getFats() *float64 {
	return dailyQuota.Fats
}

func (dailyQuota UsersDailyQuotasViewWrapper) getFatsSaturated() *float64 {
	return dailyQuota.FatsSaturated
}

func (dailyQuota UsersDailyQuotasViewWrapper) getCarbs() *float64 {
	return dailyQuota.Carbs
}

func (dailyQuota UsersDailyQuotasViewWrapper) getCarbsSugar() *float64 {
	return dailyQuota.CarbsSugar
}

func (dailyQuota UsersDailyQuotasViewWrapper) getCarbsSlowRelease() *float64 {
	return dailyQuota.CarbsSlowRelease
}

func (dailyQuota UsersDailyQuotasViewWrapper) getCarbsFastRelease() *float64 {
	return dailyQuota.CarbsFastRelease
}

func (dailyQuota UsersDailyQuotasViewWrapper) getProteins() *float64 {
	return dailyQuota.Proteins
}

func (dailyQuota UsersDailyQuotasViewWrapper) getSalt() *float64 {
	return dailyQuota.Salt
}

type DailyQuota struct {
	ID               int64
	Owner            *User
	ArchivedDateTime *time.Time
	Calories         *float64
	Fats             *float64
	FatsSaturated    *float64
	Carbs            *float64
	CarbsSugar       *float64
	CarbsSlowRelease *float64
	CarbsFastRelease *float64
	Proteins         *float64
	Salt             *float64
}

func converDailyQuota(dailyQuota DailyQuotaFromDB) *DailyQuota {
	if dailyQuota.getId() == nil {
		return nil
	} else {
		return &DailyQuota{
			ID:               *dailyQuota.getId(),
			ArchivedDateTime: dailyQuota.getArchivedDateTime(),
			Calories:         dailyQuota.getCalories(),
			Fats:             dailyQuota.getFats(),
			FatsSaturated:    dailyQuota.getFatsSaturated(),
			Carbs:            dailyQuota.getCarbs(),
			CarbsSugar:       dailyQuota.getCarbsSugar(),
			CarbsSlowRelease: dailyQuota.getCarbsSlowRelease(),
			CarbsFastRelease: dailyQuota.getCarbsFastRelease(),
			Proteins:         dailyQuota.getProteins(),
			Salt:             dailyQuota.getSalt(),
		}
	}
}

func FindDailyQuotaById(ctx context.Context, id int64) (*DailyQuota, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	row, err := queries.FindDailyQuotaById(ctx, id)
	if err != nil {
		return nil, err
	}
	dailyQuota := converDailyQuota(DailyQuotaSqlcWrapper{row.DailyQuotaSqlc})
	dailyQuota.Owner = convertUser(UserSqlcWrapper{row.UserSqlc})
	return dailyQuota, nil
}

func ListDailyQuotasForUser(ctx context.Context, id int64) ([]*DailyQuota, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.ListDailyQuotasForUser(ctx, id)
	if err != nil {
		return nil, err
	}
	result := make([]*DailyQuota, len(list))
	for i, dq := range list {
		result[i] = converDailyQuota(DailyQuotaSqlcWrapper{dq.DailyQuotaSqlc})
	}
	return result, nil
}

func FindDailyQuotaByOwnerAndCurrentDay(ctx context.Context, id int64) (*DailyQuota, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	dailyQuota, err := queries.FindDailyQuotaByOwnerAndDate(ctx, sqlc.FindDailyQuotaByOwnerAndDateParams{
		OwnerID: id,
		Date:    &now,
	})
	if err != nil {
		return nil, err
	}
	return converDailyQuota(DailyQuotaSqlcWrapper{dailyQuota.DailyQuotaSqlc}), nil
}

func FindDailyQuotaByOwnerAndDate(ctx context.Context, id int64, date time.Time) (*DailyQuota, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	dailyQuota, err := queries.FindDailyQuotaByOwnerAndDate(ctx, sqlc.FindDailyQuotaByOwnerAndDateParams{
		OwnerID: id,
		Date:    &date,
	})
	if err != nil {
		return nil, err
	}
	return converDailyQuota(DailyQuotaSqlcWrapper{dailyQuota.DailyQuotaSqlc}), nil
}

type CreateDailyQuotaRequest struct {
	OwnerID          int64    `json:"owner_id"`
	Calories         *float64 `json:"calories"`
	Fats             *float64 `json:"fats"`
	FatsSaturated    *float64 `json:"fats_saturated"`
	Carbs            *float64 `json:"carbs"`
	CarbsSugar       *float64 `json:"carbs_sugar"`
	CarbsSlowRelease *float64 `json:"carbs_slow_release"`
	CarbsFastRelease *float64 `json:"carbs_fast_release"`
	Proteins         *float64 `json:"proteins"`
	Salt             *float64 `json:"salt"`
}

func CreateDailyQuota(ctx context.Context, dailyQuota CreateDailyQuotaRequest) (*DailyQuota, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	dailyQuotaSqlc, err := queries.CreateDailyQuota(ctx, sqlc.CreateDailyQuotaParams{
		OwnerID:          dailyQuota.OwnerID,
		Calories:         dailyQuota.Calories,
		Fats:             dailyQuota.Fats,
		FatsSaturated:    dailyQuota.FatsSaturated,
		Carbs:            dailyQuota.Carbs,
		CarbsSugar:       dailyQuota.CarbsSugar,
		CarbsSlowRelease: dailyQuota.CarbsSlowRelease,
		CarbsFastRelease: dailyQuota.CarbsFastRelease,
		Proteins:         dailyQuota.Proteins,
		Salt:             dailyQuota.Salt,
	})
	if err != nil {
		return nil, err
	}
	return converDailyQuota(DailyQuotaSqlcWrapper{dailyQuotaSqlc}), nil
}

type UpdateDailyQuotaRequest struct {
	ID               int64    `json:"id"`
	OwnerID          int64    `json:"owner_id"`
	Calories         *float64 `json:"calories"`
	Fats             *float64 `json:"fats"`
	FatsSaturated    *float64 `json:"fats_saturated"`
	Carbs            *float64 `json:"carbs"`
	CarbsSugar       *float64 `json:"carbs_sugar"`
	CarbsSlowRelease *float64 `json:"carbs_slow_release"`
	CarbsFastRelease *float64 `json:"carbs_fast_release"`
	Proteins         *float64 `json:"proteins"`
	Salt             *float64 `json:"salt"`
}

func UpdateDailyQuota(ctx context.Context, dailyQuota UpdateDailyQuotaRequest) (*DailyQuota, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	dailyQuotaSqlc, err := queries.UpdateDailyQuota(ctx, sqlc.UpdateDailyQuotaParams{
		OwnerID:          dailyQuota.OwnerID,
		Calories:         dailyQuota.Calories,
		Fats:             dailyQuota.Fats,
		FatsSaturated:    dailyQuota.FatsSaturated,
		Carbs:            dailyQuota.Carbs,
		CarbsSugar:       dailyQuota.CarbsSugar,
		CarbsSlowRelease: dailyQuota.CarbsSlowRelease,
		CarbsFastRelease: dailyQuota.CarbsFastRelease,
		Proteins:         dailyQuota.Proteins,
		Salt:             dailyQuota.Salt,
		ID:               dailyQuota.ID,
	})
	if err != nil {
		return nil, err
	}
	return converDailyQuota(DailyQuotaSqlcWrapper{dailyQuotaSqlc}), nil
}

func ArchiveDailyQuota(ctx context.Context, id int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.ArchiveDailyQuota(ctx, id); err != nil {
		return err
	}
	return nil
}
