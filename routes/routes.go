package routes

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ListRequest struct {
	Name string `json:"name"`
}

type ListItems struct {
	Id   int16  `json:"item_id"`
	Name string `json:"name"`
}
type ListResponse struct {
	Name      string      `json:"name"`
	ListId    string      `json:"list_id"`
	CreatedAt string      `json:"created_at"`
	Items     []ListItems `json:"items"` // include the items
}

type Response map[string]*ListResponse

func Setup(e *echo.Echo, db *sql.DB) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})
	e.Use(middleware.Logger())

	e.GET("/lists", func(c echo.Context) error {
		l := slog.Default()
		db := c.Get("db").(*sql.DB) // retrieve it
		rows, err := db.Query(`
			select 
				l.list_id, 
				l.name as list_name, 
				l.created_at, 
				li.id AS item_id,
				li.item_name
			FROM list l
			LEFT JOIN list_items li ON li.list_id = l.list_id 
			ORDER BY l.created_at
			`)
		if err != nil {
			l.Error("", "err", err)
			return echo.NewHTTPError(http.StatusBadRequest, "provide a valid input")
		}
		defer rows.Close()

		lists := make(Response)
		for rows.Next() {
			var (
				list     ListResponse
				itemId   sql.NullInt16
				itemName sql.NullString
			)
			if err := rows.Scan(&list.ListId, &list.Name, &list.CreatedAt, &itemId, &itemName); err != nil {
				return err
			}

			lists[list.ListId] = &list

			hasItems := itemId.Valid && itemName.Valid
			if !hasItems {
				l.Info("not item found", "list", list.Name)
				continue
			}

			items := ListItems{
				Id:   itemId.Int16,
				Name: itemName.String,
			}

			if _, exists := lists[list.ListId]; exists {
				lists[list.ListId].Items = append(lists[list.ListId].Items, items)
			}
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
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			ListId:    uuid.Must(uuid.NewRandom()).String(),
			Name:      req.Name,
		}
		_, err := db.Exec(
			`INSERT INTO list (list_id, name, created_at) VALUES (?, ?, ?)`,
			response.ListId, response.Name, response.CreatedAt,
		)
		if err != nil {
			e.Logger.Fatal(err)
			return echo.NewHTTPError(http.StatusBadRequest, "cannot insert into database")
		}

		c.JSON(http.StatusCreated, response)
		return nil
	})

	e.PATCH("/lists/:list_id/items/:item_id", func(c echo.Context) error {
		_ = c.Get("db").(*sql.DB) // retrieve it
		//listId := c.Param("list_id")

		//var id int
		//err := db.QueryRow("SELECT id FROM list WHERE list_id = ?", listId).Scan(&id)
		//if err != nil {
		//	if err == sql.ErrNoRows {
		//		e.Logger.Info(err)
		//		return echo.ErrNotFound
		//	}
		//	e.Logger.Fatal(err)
		//	return echo.ErrInternalServerError
		//}

		return nil
	})

	e.DELETE("/lists/:list_id", func(c echo.Context) error {
		db := c.Get("db").(*sql.DB) // retrieve it

		listId := c.Param("list_id")
		if _, err := db.Exec("delete from list where list_id=?", listId); err != nil {
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusNoContent)
	})
}
