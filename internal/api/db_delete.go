package api

type DatabasePageDeleteInput struct {
	Database string
	PageID   string
}

// DeletePages deletes a page from a database
func (db database) DeletePage(input DatabasePageDeleteInput) error {
	// TODO: implement
	return nil
}
