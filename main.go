package main

import (
	"github.com/labstack/echo"
	"github.com/onorbit/letterite/bookdb"
	"github.com/onorbit/letterite/handler"
	"github.com/onorbit/letterite/page"
)

func main() {
	if err := bookdb.Initialize("book.sqlite3"); err != nil {
		panic(err)
	}

	if err := page.Initialize(); err != nil {
		panic(err)
	}

	e := echo.New()

	// page handlers.
	e.POST("/apis/v1/page", handler.CreatePage)
	e.GET("/apis/v1/pages/parent/:parentPageID", handler.GetPagesByParent)
	e.GET("/apis/v1/page/:id", handler.GetPage)
	e.POST("/apis/v1/page/:id", handler.UpdatePage)
	e.DELETE("/apis/v1/page/:id", handler.DeletePage)

	e.Logger.Fatal(e.Start(":10900"))
}
