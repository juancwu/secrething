// All the error codes can be found here: https://www.postgresql.org/docs/current/errcodes-appendix.html
package store

// Class 23 - Integrity Constraint Violation
const (
	PG_ERR_UNIQUE_VIOLATION = "unique_violation"
)
