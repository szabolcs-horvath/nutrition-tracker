package repository

import (
	"context"
	"github.com/szabolcs-horvath/nutrition-tracker/custom_types"
	sqlc "github.com/szabolcs-horvath/nutrition-tracker/generated"
)

type Meal struct {
	ID           int64
	Owner        *User
	Notification *Notification
	Name         string
	Time         custom_types.Time
	Quotas       map[custom_types.Quota]*float64
	Archived     bool
}

func convertMeal(meal *sqlc.Meal_sqlc) *Meal {
	quotas := map[custom_types.Quota]*float64{
		custom_types.Calories:         meal.CaloriesQuota,
		custom_types.Fats:             meal.FatsQuota,
		custom_types.FatsSaturated:    meal.FatsSaturatedQuota,
		custom_types.Carbs:            meal.CarbsQuota,
		custom_types.CarbsSugar:       meal.CarbsSugarQuota,
		custom_types.CarbsSlowRelease: meal.CarbsSlowReleaseQuota,
		custom_types.CarbsFastRelease: meal.CarbsFastReleaseQuota,
		custom_types.Proteins:         meal.ProteinsQuota,
		custom_types.Salt:             meal.SaltQuota,
	}
	return &Meal{
		ID:       meal.ID,
		Name:     meal.Name,
		Time:     meal.Time,
		Quotas:   quotas,
		Archived: meal.Archived,
	}
}

func FindMealById(ctx context.Context, id int64) (*Meal, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	row, err := queries.FindMealById(ctx, id)
	if err != nil {
		return nil, err
	}
	meal := convertMeal(&row.MealSqlc)
	meal.Owner = convertUser(UserSqlcWrapper{row.UserSqlc})
	meal.Notification = convertNotification(MealsNotificationsViewWrapper{row.MealsNotificationsView})
	return meal, nil
}

func FindMealsForUser(ctx context.Context, ownerId int64, archived bool) ([]*Meal, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.FindMealsForUser(ctx, sqlc.FindMealsForUserParams{
		OwnerID:  ownerId,
		Archived: archived,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*Meal, len(list))
	for i, m := range list {
		result[i] = convertMeal(&m.MealSqlc)
		result[i].Owner = convertUser(UserSqlcWrapper{m.UserSqlc})
		result[i].Notification = convertNotification(MealsNotificationsViewWrapper{m.MealsNotificationsView})
	}
	return result, nil
}

type CreateMealRequest struct {
	OwnerID               int64             `json:"owner_id"`
	CreateNotification    bool              `json:"create_notification"`
	Name                  string            `json:"name"`
	Time                  custom_types.Time `json:"time"`
	CaloriesQuota         *float64          `json:"calories_quota"`
	FatsQuota             *float64          `json:"fats_quota"`
	FatsSaturatedQuota    *float64          `json:"fats_saturated_quota"`
	CarbsQuota            *float64          `json:"carbs_quota"`
	CarbsSugarQuota       *float64          `json:"carbs_sugar_quota"`
	CarbsSlowReleaseQuota *float64          `json:"carbs_slow_release_quota"`
	CarbsFastReleaseQuota *float64          `json:"carbs_fast_release_quota"`
	ProteinsQuota         *float64          `json:"proteins_quota"`
	SaltQuota             *float64          `json:"salt_quota"`
}

func CreateMeal(ctx context.Context, meal CreateMealRequest) (*Meal, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var result *Meal
	err = DoInTransaction(ctx, db, func(childCtx context.Context, queries *sqlc.Queries) error {
		var notification sqlc.Notification_sqlc
		var notificationId *int64
		if meal.CreateNotification {
			notificationSqlc, notiErr := queries.CreateNotification(childCtx, sqlc.CreateNotificationParams{
				OwnerID:      meal.OwnerID,
				Time:         meal.Time,
				DelaySeconds: nil,
				DelayDate:    nil,
			})
			if notiErr != nil {
				return notiErr
			}
			notification = notificationSqlc
			notificationId = &notificationSqlc.ID
		}
		mealSqlc, mealErr := queries.CreateMeal(childCtx, sqlc.CreateMealParams{
			OwnerID:               meal.OwnerID,
			NotificationID:        notificationId,
			Name:                  meal.Name,
			Time:                  meal.Time,
			CaloriesQuota:         meal.CaloriesQuota,
			FatsQuota:             meal.FatsQuota,
			FatsSaturatedQuota:    meal.FatsSaturatedQuota,
			CarbsQuota:            meal.CarbsQuota,
			CarbsSugarQuota:       meal.CarbsSugarQuota,
			CarbsSlowReleaseQuota: meal.CarbsSlowReleaseQuota,
			CarbsFastReleaseQuota: meal.CarbsFastReleaseQuota,
			ProteinsQuota:         meal.ProteinsQuota,
			SaltQuota:             meal.SaltQuota,
		})
		if mealErr != nil {
			return err
		}
		result = convertMeal(&mealSqlc)
		if meal.CreateNotification {
			result.Notification = convertNotification(NotificationSqlcWrapper{notification})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

type UpdateMealRequest struct {
	ID                    int64             `json:"id"`
	NotificationID        *int64            `json:"notification_id"`
	Name                  string            `json:"name"`
	Time                  custom_types.Time `json:"time"`
	CaloriesQuota         *float64          `json:"calories_quota"`
	FatsQuota             *float64          `json:"fats_quota"`
	FatsSaturatedQuota    *float64          `json:"fats_saturated_quota"`
	CarbsQuota            *float64          `json:"carbs_quota"`
	CarbsSugarQuota       *float64          `json:"carbs_sugar_quota"`
	CarbsSlowReleaseQuota *float64          `json:"carbs_slow_release_quota"`
	CarbsFastReleaseQuota *float64          `json:"carbs_fast_release_quota"`
	ProteinsQuota         *float64          `json:"proteins_quota"`
	SaltQuota             *float64          `json:"salt_quota"`
}

func UpdateMeal(ctx context.Context, meal UpdateMealRequest) (*Meal, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	mealSqlc, err := queries.UpdateMeal(ctx, sqlc.UpdateMealParams{
		NotificationID:        meal.NotificationID,
		Name:                  meal.Name,
		Time:                  meal.Time,
		CaloriesQuota:         meal.CaloriesQuota,
		FatsQuota:             meal.FatsQuota,
		FatsSaturatedQuota:    meal.FatsSaturatedQuota,
		CarbsQuota:            meal.CarbsQuota,
		CarbsSugarQuota:       meal.CarbsSugarQuota,
		CarbsSlowReleaseQuota: meal.CarbsSlowReleaseQuota,
		CarbsFastReleaseQuota: meal.CarbsFastReleaseQuota,
		ProteinsQuota:         meal.ProteinsQuota,
		SaltQuota:             meal.SaltQuota,
		ID:                    meal.ID,
	})
	if err != nil {
		return nil, err
	}
	return convertMeal(&mealSqlc), nil
}

func ArchiveMeal(ctx context.Context, mealId int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.ArchiveMeal(ctx, mealId); err != nil {
		return err
	}
	return nil
}
