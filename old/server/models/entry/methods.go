package entry

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/juancwu/konbini/server/database"
	"github.com/juancwu/konbini/server/utils"
)

// create multiple personal bento entries and returns the inserted entries ids
func CreateEntries(tx *sql.Tx, bid string, keyvals []string) error {
	if bid == "" {
		return errors.New("Bento ID must not be empty.")
	}

	// create insert values
	var builder strings.Builder
	c := 2
	data := []any{bid}
	utils.Logger().Infof("len of keyvals: %d\n", len(keyvals))
	for i := 0; i < len(keyvals)-1; i += 2 {
		builder.WriteString(fmt.Sprintf("($%d, $%d, $1)", c, c+1))
		data = append(data, keyvals[i], keyvals[i+1])
		if i < len(keyvals)-2 {
			builder.WriteString(",")
		}
		c += 2
	}

	_, err := tx.Exec(
		fmt.Sprintf("INSERT INTO personal_bento_entries (name, content, personal_bento_id) VALUES %s;", builder.String()),
		data...,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetPersonalBentoEntries(bid string) ([]string, error) {
	if bid == "" {
		return nil, errors.New("Bento ID must not be empty.")
	}

	rows, err := database.DB().Query("SELECT name, content FROM personal_bento_entries WHERE personal_bento_id = $1;", bid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keyvals := []string{}
	for rows.Next() {
		var (
			key   string
			value string
		)
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		keyvals = append(keyvals, key, value)
	}

	return keyvals, nil
}
