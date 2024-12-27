package repository

import (
	"context"
	"github.com/szabolcs-horvath/nutrition-tracker/custom_types"
	sqlc "github.com/szabolcs-horvath/nutrition-tracker/generated"
	"time"
)

type DurationWrapper struct {
	Duration *time.Duration
}

func (duration DurationWrapper) Seconds() *int64 {
	if duration.Duration != nil {
		var seconds = int64(duration.Duration.Seconds())
		return &seconds
	}
	return nil
}

type NotificationFromDB interface {
	getId() *int64
	getOwnerID() *int64
	getTime() *custom_types.Time
	getDelay() *time.Duration
	getDelayDate() *custom_types.Date
}

type NotificationSqlcWrapper struct {
	sqlc.Notification_sqlc
}

func (notification NotificationSqlcWrapper) getId() *int64 {
	return &notification.ID
}

func (notification NotificationSqlcWrapper) getOwnerID() *int64 {
	return &notification.OwnerID
}

func (notification NotificationSqlcWrapper) getTime() *custom_types.Time {
	return &notification.Time
}

func (notification NotificationSqlcWrapper) getDelay() *time.Duration {
	if notification.DelaySeconds != nil {
		duration := time.Duration(*notification.DelaySeconds) * time.Second
		return &duration
	}
	return nil
}

func (notification NotificationSqlcWrapper) getDelayDate() *custom_types.Date {
	return notification.DelayDate
}

type MealsNotificationsViewWrapper struct {
	sqlc.MealsNotificationsView
}

func (notification MealsNotificationsViewWrapper) getId() *int64 {
	return notification.ID
}

func (notification MealsNotificationsViewWrapper) getOwnerID() *int64 {
	return notification.OwnerID
}

func (notification MealsNotificationsViewWrapper) getTime() *custom_types.Time {
	return notification.Time
}

func (notification MealsNotificationsViewWrapper) getDelay() *time.Duration {
	if notification.DelaySeconds != nil {
		duration := time.Duration(*notification.DelaySeconds) * time.Second
		return &duration
	}
	return nil
}

func (notification MealsNotificationsViewWrapper) getDelayDate() *custom_types.Date {
	return notification.DelayDate
}

type Notification struct {
	ID        int64
	Owner     *User
	Time      custom_types.Time
	Delay     *time.Duration
	DelayDate *custom_types.Date
}

func convertNotification(notification NotificationFromDB) *Notification {
	if notification.getId() == nil {
		return nil
	} else {
		return &Notification{
			ID:        *notification.getId(),
			Time:      *notification.getTime(),
			Delay:     notification.getDelay(),
			DelayDate: notification.getDelayDate(),
		}
	}
}

func ListNotificationsByUserId(ctx context.Context, ownerId int64) ([]*Notification, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.ListNotificationsByUserId(ctx, ownerId)
	if err != nil {
		return nil, err
	}
	var result = make([]*Notification, len(list))
	for i, n := range list {
		result[i] = convertNotification(NotificationSqlcWrapper{n.NotificationSqlc})
		result[i].Owner = convertUser(UserSqlcWrapper{n.UserSqlc})
	}
	return result, nil
}

type CreateNotificationRequest struct {
	OwnerID      int64              `json:"owner_id"`
	Time         custom_types.Time  `json:"time"`
	DelaySeconds *time.Duration     `json:"delay_seconds"`
	DelayDate    *custom_types.Date `json:"delay_date"`
}

func CreateNotification(ctx context.Context, notification CreateNotificationRequest) (*Notification, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	notificationSqlc, err := queries.CreateNotification(ctx, sqlc.CreateNotificationParams{
		OwnerID:      notification.OwnerID,
		Time:         notification.Time,
		DelaySeconds: DurationWrapper{Duration: notification.DelaySeconds}.Seconds(),
		DelayDate:    notification.DelayDate,
	})
	if err != nil {
		return nil, err
	}
	return convertNotification(NotificationSqlcWrapper{notificationSqlc}), nil
}

type UpdateNotificationRequest struct {
	ID           int64              `json:"id"`
	OwnerID      int64              `json:"owner_id"`
	Time         custom_types.Time  `json:"time"`
	DelaySeconds *time.Duration     `json:"delay_seconds"`
	DelayDate    *custom_types.Date `json:"delay_date"`
}

func UpdateNotification(ctx context.Context, notification UpdateNotificationRequest) (*Notification, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	notificationSqlc, err := queries.UpdateNotification(ctx, sqlc.UpdateNotificationParams{
		OwnerID:      notification.OwnerID,
		Time:         notification.Time,
		DelaySeconds: DurationWrapper{Duration: notification.DelaySeconds}.Seconds(),
		DelayDate:    notification.DelayDate,
		ID:           notification.ID,
	})
	if err != nil {
		return nil, err
	}
	return convertNotification(NotificationSqlcWrapper{notificationSqlc}), nil
}

func DeleteNotification(ctx context.Context, id int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.DeleteNotification(ctx, id); err != nil {
		return err
	}
	return nil
}
