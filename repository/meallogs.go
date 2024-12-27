package repository

import (
	"context"
	sqlc "github.com/szabolcs-horvath/nutrition-tracker/generated"
	"time"
)

type MealLog struct {
	ID                int64
	Meal              *Meal
	Item              *Item
	Portion           *Portion
	PortionMultiplier float64
	DateTime          time.Time
}

func (m MealLog) getMultiplier() float64 {
	return m.PortionMultiplier * m.Portion.getMultiplier()
}

func (m MealLog) GetCalories() float64 {
	return m.getMultiplier() * m.Item.CaloriesPer100
}

func (m MealLog) GetFats() float64 {
	return m.getMultiplier() * m.Item.FatsPer100
}

func (m MealLog) GetFatsSaturated() *float64 {
	var result float64
	if m.Item.FatsSaturatedPer100 != nil {
		result = m.getMultiplier() * *m.Item.FatsSaturatedPer100
	}
	return &result
}

func (m MealLog) GetCarbs() float64 {
	return m.getMultiplier() * m.Item.CarbsPer100
}

func (m MealLog) GetCarbsSugar() *float64 {
	var result float64
	if m.Item.CarbsSugarPer100 != nil {
		result = m.getMultiplier() * *m.Item.CarbsSugarPer100
	}
	return &result
}

func (m MealLog) GetCarbsSlowRelease() *float64 {
	var result float64
	if m.Item.CarbsSlowReleasePer100 != nil {
		result = m.getMultiplier() * *m.Item.CarbsSlowReleasePer100
	}
	return &result
}

func (m MealLog) GetCarbsFastRelease() *float64 {
	var result float64
	if m.Item.CarbsFastReleasePer100 != nil {
		result = m.getMultiplier() * *m.Item.CarbsFastReleasePer100
	}
	return &result
}

func (m MealLog) GetProteins() float64 {
	return m.getMultiplier() * m.Item.ProteinsPer100
}

func (m MealLog) GetSalt() *float64 {
	var result float64
	if m.Item.SaltPer100 != nil {
		result = m.getMultiplier() * *m.Item.SaltPer100
	}
	return &result
}

func convertMealLog(mealLog *sqlc.MealLog_sqlc) *MealLog {
	return &MealLog{
		ID:                mealLog.ID,
		PortionMultiplier: mealLog.PortionMultiplier,
		DateTime:          mealLog.Datetime,
	}
}

func FindMealLogById(ctx context.Context, id int64) (*MealLog, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	row, err := queries.FindMealLogById(ctx, id)
	if err != nil {
		return nil, err
	}
	mealLog := convertMealLog(&row.MealLogSqlc)
	mealLog.Meal = convertMeal(&row.MealSqlc)
	mealLog.Item = convertItem(&row.ItemSqlc)
	mealLog.Portion = convertPortion(&row.PortionSqlc)
	return mealLog, nil
}

func FindMealLogsForUserAndCurrentDay(ctx context.Context, ownerId int64) ([]*MealLog, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.FindMealLogsForUserAndDate(ctx, sqlc.FindMealLogsForUserAndDateParams{
		OwnerID: ownerId,
		Date:    time.Now(),
	})
	if err != nil {
		return nil, err
	}
	result := make([]*MealLog, len(list))
	for i, m := range list {
		result[i] = convertMealLog(&m.MealLogSqlc)
		result[i].Meal = convertMeal(&m.MealSqlc)
		result[i].Item = convertItem(&m.ItemSqlc)
		result[i].Portion = convertPortion(&m.PortionSqlc)
	}
	return result, nil
}

func FindMealLogsForUserAndDate(ctx context.Context, ownerId int64, date time.Time) ([]*MealLog, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.FindMealLogsForUserAndDate(ctx, sqlc.FindMealLogsForUserAndDateParams{
		OwnerID: ownerId,
		Date:    date,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*MealLog, len(list))
	for i, m := range list {
		result[i] = convertMealLog(&m.MealLogSqlc)
		result[i].Meal = convertMeal(&m.MealSqlc)
		result[i].Item = convertItem(&m.ItemSqlc)
		result[i].Portion = convertPortion(&m.PortionSqlc)
	}
	return result, nil
}

type CreateMealLogRequest struct {
	MealID            int64     `json:"meal_id"`
	ItemID            int64     `json:"item_id"`
	PortionID         int64     `json:"portion_id"`
	PortionMultiplier float64   `json:"portion_multiplier"`
	DateTime          time.Time `json:"date_time"`
}

func CreateMealLog(ctx context.Context, mealLog CreateMealLogRequest) (*MealLog, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	mealLogSqlc, err := queries.CreateMealLog(ctx, sqlc.CreateMealLogParams{
		MealID:            mealLog.MealID,
		ItemID:            mealLog.ItemID,
		PortionID:         mealLog.PortionID,
		PortionMultiplier: mealLog.PortionMultiplier,
		Datetime:          mealLog.DateTime,
	})
	if err != nil {
		return nil, err
	}
	return convertMealLog(&mealLogSqlc), nil
}

type UpdateMealLogRequest struct {
	ID                int64     `json:"id"`
	MealID            int64     `json:"meal_id"`
	ItemID            int64     `json:"item_id"`
	PortionID         int64     `json:"portion_id"`
	PortionMultiplier float64   `json:"portion_multiplier"`
	DateTime          time.Time `json:"date_time"`
}

func UpdateMealLog(ctx context.Context, mealLog UpdateMealLogRequest) (*MealLog, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	mealLogSqlc, err := queries.UpdateMealLog(ctx, sqlc.UpdateMealLogParams{
		MealID:            mealLog.MealID,
		ItemID:            mealLog.ItemID,
		PortionID:         mealLog.PortionID,
		PortionMultiplier: mealLog.PortionMultiplier,
		Datetime:          mealLog.DateTime,
		ID:                mealLog.ID,
	})
	if err != nil {
		return nil, err
	}
	return convertMealLog(&mealLogSqlc), nil
}

func DeleteMealLog(ctx context.Context, id int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.DeleteMealLog(ctx, id); err != nil {
		return err
	}
	return nil
}
