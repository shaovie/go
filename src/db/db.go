package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var gDB *sql.DB

func Init(host, port, user, passwd, charset, dbName string) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&", user, passwd, host, port, dbName, charset)
	gDB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	if gDB != nil {
		gDB.Close()
	}
}

func RawQuery(querySql string, args ...interface{}) (map[int]map[string]string, error) {
	rows, err := gDB.Query(querySql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	record := make(map[int]map[string]string)
	i := 0
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		row := make(map[string]string)
		for j, v := range values {
			if v != nil {
				row[columns[j]] = string(v)
			}
		}
		if len(row) > 0 {
			record[i] = row
			i++
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return record, nil
}
