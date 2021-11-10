package page

import (
	"errors"
	"strconv"

	"github.com/onorbit/letterite/bookdb"
	"github.com/onorbit/letterite/common"
)

var (
	ErrInvalidUpdateField = errors.New("invalid update field")
)

func Initialize() error {
	return nil
}

func CreatePage(parentPageID int64, subject string) (common.Page, error) {
	newPage, err := bookdb.CreatePage(parentPageID, subject)
	return newPage, err
}

func GetPagesByParent(parentPageID int64) ([]common.PageSummary, error) {
	pages, err := bookdb.GetPagesByParent(parentPageID)
	return pages, err
}

func GetPage(pageID int64) (common.Page, error) {
	page, err := bookdb.GetPage(pageID)
	return page, err
}

func UpdatePage(pageID int64, updateFieldStr, newValueStr string) error {
	switch updateFieldStr {
	case "parentId":
		parentPageID, err := strconv.ParseInt(newValueStr, 10, 64)
		if err != nil {
			return err
		}

		return bookdb.UpdatePageParent(pageID, parentPageID)
	case "order":
		order, err := strconv.ParseInt(newValueStr, 10, 64)
		if err != nil {
			return err
		}

		return bookdb.UpdatePageOrder(pageID, order)
	case "subject":
		return bookdb.UpdatePageSubject(pageID, newValueStr)
	default:
		return ErrInvalidUpdateField
	}
}

func DeletePage(pageID int64) error {
	err := bookdb.DeletePage(pageID)
	return err
}
