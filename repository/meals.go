package repository

import (
	"context"
	sqlc "github.com/szabolcs-horvath/nutrition-tracker/generated"
	"time"
)

type Meal struct {
	ID                    int64
	Owner                 *User
	Notification          *Notification
	Name                  string
	Time                  time.Time
	CaloriesQuota         *float64
	FatsQuota             *float64
	FatsSaturatedQuota    *float64
	CarbsQuota            *float64
	CarbsSugarQuota       *float64
	CarbsSlowReleaseQuota *float64
	CarbsFastReleaseQuota *float64
	ProteinsQuota         *float64
	SaltQuota             *float64
	Archived              bool
}

func convertMeal(meal *sqlc.Meal_sqlc) *Meal {
	return &Meal{
		ID:                    meal.ID,
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
		Archived:              meal.Archived,
	}
}

func FindNonArchivedMealsForUser(ctx context.Context, ownerId int64) ([]*Meal, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.FindNonArchivedMealsForUser(ctx, ownerId)
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

func CreateMeal(ctx context.Context, meal *Meal) (*Meal, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	var result *Meal
	err = DoInTransaction(ctx, db, func(queries *sqlc.Queries) error {
		notificationSqlc, notiErr := queries.CreateNotification(ctx, sqlc.CreateNotificationParams{
			OwnerID:   meal.Notification.Owner.ID,
			Time:      meal.Notification.Time,
			Delay:     meal.Notification.Delay,
			DelayDate: meal.Notification.DelayDate,
			Name:      meal.Notification.Name,
		})
		if notiErr != nil {
			return notiErr
		}
		mealSqlc, mealErr := queries.CreateMeal(ctx, sqlc.CreateMealParams{
			OwnerID:               meal.Owner.ID,
			NotificationID:        meal.Notification.ID,
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
		result.Notification = convertNotification(NotificationSqlcWrapper{notificationSqlc})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateMeal(ctx context.Context, meal *Meal) (*Meal, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	mealSqlc, err := queries.UpdateMeal(ctx, sqlc.UpdateMealParams{
		OwnerID:               meal.Owner.ID,
		NotificationID:        meal.Notification.ID,
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
