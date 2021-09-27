package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/onorbit/letterite/page"
)

// Response structures
type CreatePageResponse struct {
	PageID       int64  `json:"pageID"`
	Subject      string `json:"subject"`
	ParentPageID int64  `json:"parentPageID"`
}

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

	pageID, err := page.CreatePage(parentPageID, subject)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	response := CreatePageResponse{
		PageID:       pageID,
		Subject:      subject,
		ParentPageID: parentPageID,
	}

	return c.JSON(http.StatusOK, response)
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
