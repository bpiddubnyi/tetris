// +build pprof

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		log.Println("start pprof web server")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
