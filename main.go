package main

import (
	"database/sql"
	"fmt"

	"github.com/1garo/easylist/env"
	"github.com/1garo/easylist/routes"
	"github.com/labstack/echo/v4"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	e := echo.New()
	env := env.Load()

	url := fmt.Sprintf("%s?authToken=%s", env.DatabaseURL, env.DatabaseToken)
	db, err := sql.Open("libsql", url)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	routes.Setup(e, db)

	e.Logger.Info("starting server on port 8000...")
	e.Logger.Fatal(e.Start(":8000"))
}
