package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"ilog"
)

var gDB *sql.DB

func Init(host, port, user, passwd, charset, dbName string) (err error) {
	gDB, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&", user, passwd, host, port, dbName, charset))
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
func RawQuery(querySql string) (map[int]map[string]string, error) {
	rows, err := gDB.Query(querySql)
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

func mysqlTest(ch chan<- int) error {
	f := func() {
		ch <- 1
	}
	defer f()
	err := Init("127.0.0.1", "3306", "root", "taojinzi", "utf8", "wemall")
	if err != nil {
		ilog.Error("mysql - " + err.Error())
		return err
	}
	//defer Close()

	n := 0
	for i := 0; i < 100000; i++ {
		rows, err := RawQuery("select * from ws_group_buy")
		if err != nil {
			ilog.Error("mysql - " + err.Error())
			return err
		}
		n += len(rows)
		//ilog.Rinfo("count = %d content=%s", len(rows), rows[0]["content"])
	}
	ilog.Rinfo("count = %d", n)

	return nil
}

func TestMyQuery() {
	closeCh := make(chan int)
	for i := 0; i < 3; i++ {
		go mysqlTest(closeCh)
	}
	for i := 0; i < 3; i++ {
		<-closeCh
	}
	os.Exit(0)
}
