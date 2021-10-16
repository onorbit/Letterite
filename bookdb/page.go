package bookdb

import (
	"errors"

	"github.com/onorbit/letterite/common"
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

func CreatePage(parentPageID int64, subject string) (newPage common.Page, err error) {
	newPage = common.Page{
		ID:           common.InvalidPageID,
		ParentPageID: parentPageID,
		Subject:      subject,
	}

	tx, err := gDatabase.Begin()
	if err != nil {
		return newPage, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			newPage.ID = common.InvalidPageID
			err = r.(error)
		}
	}()

	// if the page belongs to existing parent, perform check.
	if parentPageID != common.RootPageID {
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
	newPageID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	newPage.ID = newPageID
	return
}

func GetPagesByParent(parentPageID int64) ([]common.PageSummary, error) {
	type PageSummary struct {
		ID                int64  `db:"id"`
		ParentPageID      int64  `db:"parent_id"`
		Subject           string `db:"subject"`
		ChildrenPageCount int    `db:"children_count"`
		ContentCount      int    `db:"content_count"`
	}

	query := `
		SELECT A.*, COUNT(B.id) AS children_count, COUNT(C.page_id) AS content_count
		FROM pages A
			LEFT JOIN pages B ON B.parent_id = A.id
			LEFT JOIN page_contents C ON C.page_id = A.id
		WHERE A.parent_id = $1
		GROUP BY A.id`

	entries := []PageSummary{}
	if err := gDatabase.Select(&entries, query, parentPageID); err != nil {
		return nil, err
	}

	ret := make([]common.PageSummary, len(entries))
	for i, entry := range entries {
		ret[i].ID = entry.ID
		ret[i].ParentPageID = entry.ParentPageID
		ret[i].Subject = entry.Subject
		ret[i].ChildrenPageCount = entry.ChildrenPageCount
		ret[i].ContentCount = entry.ContentCount
	}

	return ret, nil
}
