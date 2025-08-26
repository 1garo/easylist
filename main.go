package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	_ "github.com/tursodatabase/go-libsql"
)

type ListRequest struct {
	Name string `json:"name"`
}

type ListResponse struct {
	Name      string    `json:"name"`
	ItemId    string    `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	e := echo.New()
	dbName := "file:./list.db"

	db, err := sql.Open("libsql", dbName)
	if err != nil {
		e.Logger.Fatalf("failed to open db %s", err)
	}
	defer db.Close()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.GET("/lists", func(c echo.Context) error {
		db := c.Get("db").(*sql.DB) // retrieve it
		rows, err := db.Query("select item_id, name, created_at from list")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "provide a valid input")
		}
		defer rows.Close()

		var lists []ListResponse
		for rows.Next() {
			var list ListResponse
			if err := rows.Scan(&list.ItemId, &list.Name, &list.CreatedAt); err != nil {
				return err
			}
			lists = append(lists, list)
		}

		return c.JSON(http.StatusOK, lists)
	})

	e.POST("/list", func(c echo.Context) error {
		req := new(ListRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "provide a valid input")
		}

		db := c.Get("db").(*sql.DB) // retrieve it
		response := ListResponse{
			CreatedAt: time.Now(),
			ItemId:    uuid.Must(uuid.NewRandom()).String(),
			Name:      req.Name,
		}
		_, err := db.Exec(
			`INSERT INTO list (item_id, name, created_at) VALUES (?, ?, ?)`,
			response.ItemId, response.Name, response.CreatedAt,
		)
		if err != nil {
			e.Logger.Fatal(err)
			return echo.NewHTTPError(http.StatusBadRequest, "cannot insert into database")
		}
		c.JSON(http.StatusOK, response)
		return nil
	})

	e.Logger.Info("starting server on port 8000...")
	e.Logger.Fatal(e.Start(":8000"))
}
