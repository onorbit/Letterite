package bookdb

func initTags() error {
	// table schemas
	schemaTags := `
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY,
			tag TEXT
		)`

	schemaPagesTags := `
		CREATE TABLE IF NOT EXISTS pages_tags (
			page_id INTEGER,
			tag_id INTEGER
		)`

	// prepare tags table and related index
	if _, err := gDatabase.Exec(schemaTags); err != nil {
		return err
	}

	if _, err := gDatabase.Exec("CREATE INDEX IF NOT EXISTS tags__tag ON tags (tag)"); err != nil {
		return err
	}

	// prepare pages-tags relation table and related indices
	if _, err := gDatabase.Exec(schemaPagesTags); err != nil {
		return err
	}

	if _, err := gDatabase.Exec("CREATE INDEX IF NOT EXISTS pages_tags__page_id ON pages_tags (page_id)"); err != nil {
		return err
	}

	if _, err := gDatabase.Exec("CREATE INDEX IF NOT EXISTS pages_tags__tag_id ON pages_tags (tag_id)"); err != nil {
		return err
	}

	return nil
}
