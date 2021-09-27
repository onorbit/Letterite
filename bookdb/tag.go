package bookdb

func initTags() error {
	// create tags table
	stmt, err := gDatabase.Prepare(`CREATE TABLE IF NOT EXISTS tags (
										id INTEGER PRIMARY KEY,
										tag TEXT
									)`)
	if err != nil {
		return err
	}
	stmt.Exec()

	// create tags-pages relation table
	stmt, err = gDatabase.Prepare(`CREATE TABLE IF NOT EXISTS pages_tags (
										page_id INTEGER,
										tag_id INTEGER
									)`)
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = gDatabase.Prepare("CREATE INDEX IF NOT EXISTS pages_tags__page_id ON pages_tags (page_id)")
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = gDatabase.Prepare("CREATE INDEX IF NOT EXISTS pages_tags__tag_id ON pages_tags (tag_id)")
	if err != nil {
		return err
	}
	stmt.Exec()

	return nil
}
