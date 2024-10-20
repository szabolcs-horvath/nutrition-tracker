package http_server

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"shorvath/nutrition-tracker/generated"
	"strconv"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"test\":\"value\"}"))
	if err != nil {
		log.Fatal(err)
	}
}

func HandleItem(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "sqlite/nutrition-tracker.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	queries := repository.New(db)

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		idParam := r.Form.Get("id")
		if idParam != "" {
			id, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			item, err := queries.FindItemById(ctx, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			jsonResponse, err := json.Marshal(item)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = w.Write(jsonResponse)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			return
		}
		http.Error(w, "You need to specify the id of the item!", http.StatusBadRequest)
		return
	case "POST":
		tx, err := db.Begin()
		defer tx.Rollback()
		var item repository.Item
		err = parseJson(r.Body, &item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//nutrition, err := queries.CreateNutrition(ctx, repository.CreateNutritionParams{
		//	CaloriesPer100g:      ,
		//	FatsPer100g:          0,
		//	FatsSaturatedPer100g: sql.NullFloat64{},
		//	CarbsPer100g:         0,
		//	CarbsSugarPer100g:    sql.NullFloat64{},
		//	ProteinsPer100g:      0,
		//	SaltPer100g:          sql.NullFloat64{},
		//})
		//item, err := queries.CreateItem(ctx, repository.CreateItemParams{
		//	Name:        name,
		//	NutritionID: nutrition.ID,
		//	Icon:        []byte(icon),
		//})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		jsonResponse, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		err = tx.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func parseJson(requestBody io.Reader, item any) error {
	return json.NewDecoder(requestBody).Decode(&item)
}
