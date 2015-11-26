package db

import (
	"fmt"
	"testing"
)

func mysqlTest(ch chan<- int) error {
	f := func() {
		ch <- 1
	}
	defer f()
	err := Init("127.0.0.1", "3306", "root", "taojinzi", "utf8", "wemall")
	if err != nil {
		fmt.Println("mysql - " + err.Error())
		return err
	}
	//defer Close()

	n := 0
	for i := 0; i < 100000; i++ {
		rows, err := RawQuery("select * from ws_group_buy where id=?", 9)
		if err != nil {
			fmt.Println("mysql - " + err.Error())
			return err
		}
		n += len(rows)
		//ilog.Rinfo("count = %d content=%s", len(rows), rows[0]["content"])
	}
	fmt.Printf("count = %d\n", n)

	return nil
}

func BenchmarkQuery(b *testing.B) {
	closeCh := make(chan int)
	for i := 0; i < 3; i++ {
		go mysqlTest(closeCh)
	}
	for i := 0; i < 3; i++ {
		<-closeCh
	}
}
