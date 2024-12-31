package repository

import (
	"context"
	"github.com/szabolcs-horvath/nutrition-tracker/custom_types"
	sqlc "github.com/szabolcs-horvath/nutrition-tracker/generated"
	"reflect"
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
	return m.PortionMultiplier * m.Portion.getUnitPerPortion() / 100
}

func (m MealLog) GetByQuota(quota custom_types.Quota) float64 {
	itemValue := reflect.ValueOf(m.Item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}
	field := itemValue.FieldByName(quota.String() + "Per100")
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return 0
		} else {
			return m.getMultiplier() * field.Elem().Interface().(float64)
		}
	} else {
		return m.getMultiplier() * field.Interface().(float64)
	}
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

func FindMealLogsForMealAndCurrentDay(ctx context.Context, mealId int64) ([]*MealLog, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.FindMealLogsForMealAndDate(ctx, sqlc.FindMealLogsForMealAndDateParams{
		MealID: mealId,
		Date:   time.Now(),
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
