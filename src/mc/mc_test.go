package mc

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func mcTest(ch chan<- int, key string) error {
	err := Del("hello" + key)
	if err != nil {
		fmt.Println("redis - " + err.Error())
		ch <- 1
		return err
	}
	for i := 0; i < 1000000; i++ {
		err = Incr("hello" + key)
		if err != nil {
			fmt.Println("redis - " + err.Error())
			ch <- 1
			return err
		}
	}
	ret, _ := Get("hello" + key)
	fmt.Println("redis - " + fmt.Sprintf("%s", ret))
	ch <- 1
	return nil
}

func BenchmarkMc(b *testing.B) {
	err := Init("172.18.8.24", "6379", "xx")
	if err != nil {
		fmt.Println("redis - " + err.Error())
		os.Exit(0)
	}

	closeCh := make(chan int)
	for i := 0; i < 3; i++ {
		go mcTest(closeCh, strconv.Itoa(i))
	}
	for i := 0; i < 3; i++ {
		<-closeCh
	}
}
