package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DBConn struct {
    db     *DB
}

func NewDB(host, port, user, passwd, charset, name string) (*DBConn, error) {
    db, err := sql.Open("mysql",
        fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&", user, passwd, host, port, name, charset))
	if err != nil {
		return fmt.Errorf("mysql - open fail[%s]", err.Error())
	}
    return &DBConn {
        db: db,
    }
}

func (c *DBConn) RawQuery(sql string) (map[{}interface]{}interface, error) {
    rows, err := c.db.Query(sql)
    if err != nil {
        return nil, fmt.Errorf("mysql - query error [%s]", err.Error())
    }

    defer rows.Close()

func handleResult()
    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        return nil, fmt.Errorf("mysql - columns error [%s]", err.Error())
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
    for rows.Next() {
        // get RawBytes from data
        err = rows.Scan(scanArgs...)
        if err != nil {
            ilog.Error("mysql - " + err.Error())
            rows.Close()
            return err
        }
        record := make(map[string]string)
        for i, col := range values {
            if col != nil {
                record[columns[i]] = string(col)
            }
        }
    }
    if err = rows.Err(); err != nil {
        ilog.Error("mysql - " + err.Error())
        rows.Close()
        return err
    }
    rows.Close()
}
}
func mysqlTest(ch chan<- int) error {
	f := func() {
		ch <- 1
	}
	defer f()

	db, err := sql.Open("mysql", "root:taojinzi@tcp(127.0.0.1:3306)/wemall?charset=utf8")
	if err != nil {
		ilog.Error("mysql - " + err.Error())
		return err
	}
	for i := 0; i < 100000; i++ {
		rows, err := db.Query("select * from ws_group_buy")
		if err != nil {
			ilog.Error("mysql - " + err.Error())
			return err
		}

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			ilog.Error("mysql - " + err.Error())
            rows.Close()
			return err
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
		for rows.Next() {
			// get RawBytes from data
			err = rows.Scan(scanArgs...)
			if err != nil {
				ilog.Error("mysql - " + err.Error())
                rows.Close()
				return err
			}
			record := make(map[string]string)
			for i, col := range values {
				if col != nil {
					record[columns[i]] = string(col)
				}
			}
		}
		if err = rows.Err(); err != nil {
			ilog.Error("mysql - " + err.Error())
            rows.Close()
			return err
		}
        rows.Close()
	}
	ilog.Rinfo("count = " + strconv.Itoa(count))

	defer db.Close()
	return nil
}


func doMysqlTest() {
	closeCh := make(chan int)
	for i := 0; i < 3; i++ {
		go mysqlTest(closeCh)
	}
	for i := 0; i < 3; i++ {
		<-closeCh
	}
	os.Exit(0)
}
