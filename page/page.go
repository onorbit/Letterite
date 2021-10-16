package page

import (
	"github.com/onorbit/letterite/bookdb"
	"github.com/onorbit/letterite/common"
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
