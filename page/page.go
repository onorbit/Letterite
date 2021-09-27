package page

import "github.com/onorbit/letterite/bookdb"

func Initialize() error {
	return nil
}

func CreatePage(parentPageID int64, subject string) (int64, error) {
	newPageID, err := bookdb.CreatePage(parentPageID, subject)
	return newPageID, err
}
