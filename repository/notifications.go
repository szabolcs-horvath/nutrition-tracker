package repository

import (
	"context"
	sqlc "github.com/szabolcs-horvath/nutrition-tracker/generated"
	"time"
)

type NotificationFromDB interface {
	getId() *int64
	getOwnerID() *int64
	getTime() *time.Time
	getDelay() *time.Time
	getDelayDate() *time.Time
	getName() *string
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

func (notification NotificationSqlcWrapper) getTime() *time.Time {
	return &notification.Time
}

func (notification NotificationSqlcWrapper) getDelay() *time.Time {
	return notification.Delay
}

func (notification NotificationSqlcWrapper) getDelayDate() *time.Time {
	return notification.DelayDate
}

func (notification NotificationSqlcWrapper) getName() *string {
	return &notification.Name
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

func (notification MealsNotificationsViewWrapper) getTime() *time.Time {
	return notification.Time
}

func (notification MealsNotificationsViewWrapper) getDelay() *time.Time {
	return notification.Delay
}

func (notification MealsNotificationsViewWrapper) getDelayDate() *time.Time {
	return notification.DelayDate
}

func (notification MealsNotificationsViewWrapper) getName() *string {
	return notification.Name
}

type Notification struct {
	ID        int64
	Owner     *User
	Time      time.Time
	Delay     *time.Time
	DelayDate *time.Time
	Name      string
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
			Name:      *notification.getName(),
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

func CreateNotification(ctx context.Context, notification *Notification) (*Notification, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	notificationSqlc, err := queries.CreateNotification(ctx, sqlc.CreateNotificationParams{
		OwnerID:   notification.Owner.ID,
		Time:      notification.Time,
		Delay:     notification.Delay,
		DelayDate: notification.DelayDate,
		Name:      notification.Name,
	})
	if err != nil {
		return nil, err
	}
	return convertNotification(NotificationSqlcWrapper{notificationSqlc}), nil
}

func UpdateNotification(ctx context.Context, notification *Notification) (*Notification, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	notificationSqlc, err := queries.UpdateNotification(ctx, sqlc.UpdateNotificationParams{
		OwnerID:   notification.Owner.ID,
		Time:      notification.Time,
		Delay:     notification.Delay,
		DelayDate: notification.DelayDate,
		Name:      notification.Name,
		ID:        notification.ID,
	})
	if err != nil {
		return nil, err
	}
	return convertNotification(NotificationSqlcWrapper{notificationSqlc}), nil
}

func DeleteNotification(ctx context.Context, notificationId int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.DeleteNotification(ctx, notificationId); err != nil {
		return err
	}
	return nil
}
