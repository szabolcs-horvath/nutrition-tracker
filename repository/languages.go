package repository

import "github.com/szabolcs-horvath/nutrition-tracker/generated"

type LanguageFromDB interface {
	getId() *int64
	getName() *string
	getNativeName() *string
}

type LanguageSqlcWrapper struct {
	sqlc.Language_sqlc
}

func (language LanguageSqlcWrapper) getId() *int64 {
	return &language.ID
}

func (language LanguageSqlcWrapper) getName() *string {
	return &language.Name
}

func (language LanguageSqlcWrapper) getNativeName() *string {
	return &language.NativeName
}

type PortionsLanguagesViewWrapper struct {
	sqlc.PortionsLanguagesView
}

func (language PortionsLanguagesViewWrapper) getId() *int64 {
	return language.ID
}

func (language PortionsLanguagesViewWrapper) getName() *string {
	return language.Name
}

func (language PortionsLanguagesViewWrapper) getNativeName() *string {
	return language.NativeName
}

type Language struct {
	ID         int64
	Name       string
	NativeName string
}

func convertLanguage(language LanguageFromDB) *Language {
	if language.getId() == nil {
		return nil
	} else {
		return &Language{
			ID:         *language.getId(),
			Name:       *language.getName(), //TODO make these nil-safe by creating some util function
			NativeName: *language.getNativeName(),
		}
	}
}
