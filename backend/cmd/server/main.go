package main

import (
	"github.com/juancwu/secrething/db"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	conn, err := db.Connect()
	if err != nil {
		e.Logger.Fatal(err)
	}

	conn.Ping()

	if err := e.Start(":3000"); err != nil {
		e.Logger.Fatal(err)
	}
}
