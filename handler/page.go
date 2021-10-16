package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/onorbit/letterite/page"
)

// CreatePage makes a new blank page with given subject and parentID.
func CreatePage(c echo.Context) error {
	subject := c.FormValue("subject")
	parentPageIDStr := c.FormValue("parentPageID")

	if len(subject) == 0 || len(parentPageIDStr) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	parentPageID, err := strconv.ParseInt(parentPageIDStr, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	page, err := page.CreatePage(parentPageID, subject)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, page)
}

// GetPagesByParent returns brief information of children pages under designated parent.
func GetPagesByParent(c echo.Context) error {
	parentPageIDStr := c.Param("parentPageID")
	if len(parentPageIDStr) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	parentPageId, err := strconv.ParseInt(parentPageIDStr, 16, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	pages, err := page.GetPagesByParent(parentPageId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pages)
}

func GetPage(c echo.Context) error {
	// TODO : implement this.
	return nil
}

func UpdatePage(c echo.Context) error {
	// TODO : implement this.
	return nil
}

func DeletePage(c echo.Context) error {
	// TODO : implement this.
	return nil
}
