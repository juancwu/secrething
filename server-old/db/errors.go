package db

// IsNoRows checks whether the given error is a no rows returned error
func IsNoRows(err error) bool {
	return err.Error() == "sql: no rows in result set"
}
