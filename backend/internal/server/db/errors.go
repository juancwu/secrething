package db

// IsNoRows checks whether the given error is a no rows returned error
func IsNoRows(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "sql: no rows in result set"
}
