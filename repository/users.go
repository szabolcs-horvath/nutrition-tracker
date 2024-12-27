package repository

import (
	"context"
	"github.com/szabolcs-horvath/nutrition-tracker/generated"
)

type UserFromDB interface {
	getId() *int64
	getLanguageID() *int64
	getDailyQuotaID() *int64
}

type UserSqlcWrapper struct {
	sqlc.User_sqlc
}

func (user UserSqlcWrapper) getId() *int64 {
	return &user.ID
}

func (user UserSqlcWrapper) getLanguageID() *int64 {
	return &user.LanguageID
}

func (user UserSqlcWrapper) getDailyQuotaID() *int64 {
	return user.DailyQuotaID
}

type PortionsUsersViewWrapper struct {
	sqlc.PortionsUsersView
}

func (user PortionsUsersViewWrapper) getId() *int64 {
	return user.ID
}

func (user PortionsUsersViewWrapper) getLanguageID() *int64 {
	return user.LanguageID
}

func (user PortionsUsersViewWrapper) getDailyQuotaID() *int64 {
	return user.DailyQuotaID
}

type ItemsUsersViewWrapper struct {
	sqlc.ItemsUsersView
}

func (user ItemsUsersViewWrapper) getId() *int64 {
	return user.ID
}

func (user ItemsUsersViewWrapper) getLanguageID() *int64 {
	return user.LanguageID
}

func (user ItemsUsersViewWrapper) getDailyQuotaID() *int64 {
	return user.DailyQuotaID
}

type User struct {
	ID         int64
	Language   *Language
	DailyQuota *DailyQuota
}

func convertUser(user UserFromDB) *User {
	if user.getId() == nil {
		return nil
	} else {
		return &User{
			ID: *user.getId(),
		}
	}
}

func ListUsers(ctx context.Context) ([]*User, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	list, err := queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*User, len(list))
	for i, u := range list {
		result[i] = convertUser(UserSqlcWrapper{u.UserSqlc})
		result[i].Language = convertLanguage(LanguageSqlcWrapper{u.LanguageSqlc})
		result[i].DailyQuota = converDailyQuota(UsersDailyQuotasViewWrapper{u.UsersDailyQuotasView})
	}
	return result, nil
}

func FindUserById(ctx context.Context, id int64) (*User, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	row, err := queries.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	user := convertUser(UserSqlcWrapper{row.UserSqlc})
	user.Language = convertLanguage(LanguageSqlcWrapper{row.LanguageSqlc})
	user.DailyQuota = converDailyQuota(UsersDailyQuotasViewWrapper{row.UsersDailyQuotasView})
	return user, nil
}

type CreateUserRequest struct {
	LanguageID   int64  `json:"language_id"`
	DailyQuotaID *int64 `json:"daily_quota_id"`
}

func CreateUser(ctx context.Context, user CreateUserRequest) (*User, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	userSqlc, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		LanguageID:   user.LanguageID,
		DailyQuotaID: user.DailyQuotaID,
	})
	if err != nil {
		return nil, err
	}
	return convertUser(UserSqlcWrapper{userSqlc}), nil
}

type UpdateUserRequest struct {
	ID           int64  `json:"id"`
	LanguageID   int64  `json:"language_id"`
	DailyQuotaID *int64 `json:"daily_quota_id"`
}

func UpdateUser(ctx context.Context, user UpdateUserRequest) (*User, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	userSqlc, err := queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		LanguageID:   user.LanguageID,
		ID:           user.ID,
		DailyQuotaID: user.DailyQuotaID,
	})
	if err != nil {
		return nil, err
	}
	return convertUser(UserSqlcWrapper{userSqlc}), nil
}

func DeleteUser(ctx context.Context, id int64) error {
	queries, err := GetQueries()
	if err != nil {
		return err
	}
	if err = queries.DeleteItem(ctx, id); err != nil {
		return err
	}
	return nil
}
