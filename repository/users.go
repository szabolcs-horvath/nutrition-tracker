package repository

import (
	"context"
	"github.com/szabolcs-horvath/nutrition-tracker/generated"
)

type UserFromDB interface {
	getId() *int64
}

type UserSqlcWrapper struct {
	sqlc.User_sqlc
}

func (user UserSqlcWrapper) getId() *int64 {
	return &user.ID
}

type PortionsUsersViewWrapper struct {
	sqlc.PortionsUsersView
}

func (user PortionsUsersViewWrapper) getId() *int64 {
	return user.ID
}

type ItemsUsersViewWrapper struct {
	sqlc.ItemsUsersView
}

func (user ItemsUsersViewWrapper) getId() *int64 {
	return user.ID
}

type User struct {
	ID       int64
	Language *Language
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
	return user, nil
}

func CreateUser(ctx context.Context, languageId int64) (*User, error) {
	queries, err := GetQueries()
	if err != nil {
		return nil, err
	}
	userSqlc, err := queries.CreateUser(ctx, languageId)
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
