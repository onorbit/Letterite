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

	parentPageID, err := strconv.ParseInt(parentPageIDStr, 16, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// TODO : in case of not found?
	pages, err := page.GetPagesByParent(parentPageID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pages)
}

func GetPage(c echo.Context) error {
	pageIDStr := c.Param("pageID")
	if len(pageIDStr) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// TODO : in case of not found?
	page, err := page.GetPage(pageID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, page)
}

func UpdatePage(c echo.Context) error {
	pageIDStr := c.Param("pageID")
	if len(pageIDStr) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	updateFieldStr := c.FormValue("updateField")
	newValueStr := c.FormValue("newValue")

	if len(updateFieldStr) == 0 || len(newValueStr) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = page.UpdatePage(pageID, updateFieldStr, newValueStr)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func DeletePage(c echo.Context) error {
	pageIDStr := c.Param("pageID")
	if len(pageIDStr) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// TODO : in case of not found or not in recycle bin?
	err = page.DeletePage(pageID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
