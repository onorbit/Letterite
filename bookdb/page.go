package bookdb

import (
	"database/sql"
	"errors"

	"github.com/onorbit/letterite/consts"
)

var (
	gStmtPageExists *sql.Stmt
	gStmtInsertPage *sql.Stmt

	ErrParentPageNotFound = errors.New("parent page not found")
)

func initPages() error {
	// create pages table
	stmt, err := gDatabase.Prepare(`CREATE TABLE IF NOT EXISTS pages (
										id INTEGER PRIMARY KEY,
										parent_id INTEGER,
										subject TEXT
									)`)
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = gDatabase.Prepare("CREATE INDEX IF NOT EXISTS pages__parent_id ON pages (parent_id)")
	if err != nil {
		return err
	}
	stmt.Exec()

	// create page_contents table
	stmt, err = gDatabase.Prepare(`CREATE TABLE IF NOT EXISTS page_contents (
										page_id INTEGER,
										revision INTEGER,
										content TEXT,
										committed_time INTEGER,
										PRIMARY KEY (page_id, revision DESC)
									) WITHOUT ROWID`)
	if err != nil {
		return err
	}
	stmt.Exec()

	// utility function for preparing statements
	prepareStmt := func(sql string, targetStmt **sql.Stmt) error {
		stmt, err = gDatabase.Prepare(sql)
		if err != nil {
			return err
		}

		*targetStmt = stmt
		return nil
	}

	// prepare statements
	if err = prepareStmt(`SELECT 1
							FROM pages
							WHERE id = ?`,
		&gStmtPageExists); err != nil {
		return err
	}

	if err = prepareStmt(`INSERT INTO pages (parent_id, subject)
							VALUES (?, ?)`,
		&gStmtInsertPage); err != nil {
		return err
	}

	return nil
}

func CreatePage(parentPageID int64, subject string) (int64, error) {
	tx, err := gDatabase.Begin()

	// if the page belongs to existing parent, perform check.
	if parentPageID != consts.InvalidPageID {
		stmt := tx.Stmt(gStmtPageExists)
		rows, err := stmt.Query(parentPageID)
		if err != nil {
			stmt.Close()
			tx.Rollback()

			return consts.InvalidPageID, err
		}

		// parent page not found.
		if rows.Next() != true {
			stmt.Close()
			tx.Rollback()

			return consts.InvalidPageID, ErrParentPageNotFound
		}

		stmt.Close()
	}

	// insert the page.
	stmt := tx.Stmt(gStmtInsertPage)
	result, err := stmt.Exec(parentPageID, subject)
	defer stmt.Close()

	if err != nil {
		tx.Rollback()
		return consts.InvalidPageID, err
	}

	newPageID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return consts.InvalidPageID, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return consts.InvalidPageID, err
	}

	return newPageID, nil
}
