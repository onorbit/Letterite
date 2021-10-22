package bookdb

import (
	"errors"

	"github.com/onorbit/letterite/common"
)

var (
	ErrPageNotFound          = errors.New("page not found")
	ErrParentPageNotFound    = errors.New("parent page not found")
	ErrPageIsNotInRecycleBin = errors.New("page is not in recycle bin")
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

func isPageExists(pageID int64) (bool, error) {
	rows, err := gDatabase.Query("SELECT 1 FROM pages WHERE id = ?", pageID)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	return rows.Next(), nil
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
		if isParentExists, err := isPageExists(parentPageID); err != nil {
			panic(err)
		} else if !isParentExists {
			panic(ErrParentPageNotFound)
		}
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

func UpdatePage(pageID, newParentPageID int64, newSubject string) (err error) {
	tx, err := gDatabase.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = r.(error)
		}
	}()

	// check if the page exists.
	if isPageExists, err := isPageExists(pageID); err != nil {
		panic(err)
	} else if !isPageExists {
		panic(ErrPageNotFound)
	}

	// if the new parent is existing page, check it.
	if newParentPageID != common.RootPageID && newParentPageID != common.RecycleBinPageID {
		if isPageExists, err := isPageExists(newParentPageID); err != nil {
			panic(err)
		} else if !isPageExists {
			panic(ErrParentPageNotFound)
		}
	}

	// perform the update.
	gDatabase.MustExec("UPDATE pages SET parent_id = ?, subject = ? WHERE id = ?", newParentPageID, newSubject, pageID)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return nil
}

func DeletePage(pageID int64) (err error) {
	tx, err := gDatabase.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = r.(error)
		}
	}()

	// the page should be in recycle bin.
	cursorPageID := pageID
	for {
		if err := gDatabase.Get(cursorPageID, "SELECT parent_id FROM pages WHERE id = ?", cursorPageID); err != nil {
			panic(err)
		}

		if cursorPageID == common.RootPageID {
			panic(ErrPageIsNotInRecycleBin)
		} else if cursorPageID == common.RecycleBinPageID {
			break
		}
	}

	// perform deletion.
	result, err := gDatabase.Exec("DELETE FROM pages WHERE id = ?", pageID)
	if err != nil {
		panic(err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		panic(err)
	} else if rowsAffected != 1 {
		panic(ErrPageNotFound)
	}

	return nil
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
