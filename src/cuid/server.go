package main

import (
	"fmt"
	"io"
	"net/http"

	"ilog"
	"mc"
)

func startServer() error {
	index := func(w http.ResponseWriter, r *http.Request) {
		err := mc.Set("PING", "x")
		if err != nil {
			ilog.Error("redis - " + err.Error())
		} else {
			io.WriteString(w, fmt.Sprintf("%s", "ok"))
		}
	}
	http.HandleFunc("/", index)
	http.ListenAndServe(":8081", nil)

	return nil
}

/*
	listener, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		ilog.Error("listen fail " + err.Error())
		os.Exit(0)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			ilog.Error("accept error " + err.Error())
			break
		}
		go func(conn net.Conn) {
			for {
				buf := make([]byte, 1024)
				_, err := conn.Read(buf)
				if err != nil {
					conn.Close()
					return
				}

				_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nConnection: keep-alive\r\nContent-Length: 5\r\n\r\nhello"))
			}
		}(conn)
	}
*/
