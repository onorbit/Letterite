package bookdb

import (
	"errors"

	"github.com/onorbit/letterite/consts"
)

var (
	ErrParentPageNotFound = errors.New("parent page not found")
)

func initPages() error {
	// table schemas
	schemaPages := `
		CREATE TABLE IF NOT EXISTS pages (
			id INTEGER PRIMARY KEY,
			parent_id INTEGER,
			subject TEXT
		)`

	schemaPageContents := `
		CREATE TABLE IF NOT EXISTS page_contents (
			page_id INTEGER,
			revision INTEGER,
			content TEXT,
			committed_time INTEGER,
			PRIMARY KEY (page_id, revision DESC)
		) WITHOUT ROWID`

	// prepare pages table and related index
	if _, err := gDatabase.Exec(schemaPages); err != nil {
		return err
	}

	if _, err := gDatabase.Exec("CREATE INDEX IF NOT EXISTS pages__parent_id ON pages (parent_id)"); err != nil {
		return err
	}

	// prepare page_contents table
	if _, err := gDatabase.Exec(schemaPageContents); err != nil {
		return err
	}

	return nil
}

func CreatePage(parentPageID int64, subject string) (newPageID int64, err error) {
	tx, err := gDatabase.Begin()
	if err != nil {
		return consts.InvalidPageID, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			newPageID = consts.InvalidPageID
			err = r.(error)
		}
	}()

	// if the page belongs to existing parent, perform check.
	if parentPageID != consts.RootPageID {
		rows, err := gDatabase.Query("SELECT 1 FROM pages WHERE id = ?", parentPageID)
		if err != nil {
			rows.Close()
			panic(err)
		}

		if !rows.Next() {
			rows.Close()
			panic(ErrParentPageNotFound)
		}

		rows.Close()
	}

	// insert the page.
	result := gDatabase.MustExec("INSERT INTO pages (parent_id, subject) VALUES (?, ?)", parentPageID, subject)
	newPageID, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return
}
