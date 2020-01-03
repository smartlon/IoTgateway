package main

import (
	"time"
)

const (
	SERVICEURL = "http://202.117.43.212:8081"
)

var (
	seed string
	sideKey string
	containerListen  map[string]bool
)

func main() {
	stop := make(map[string] chan string)
	containerListen = make(map[string]bool)
	for {
		results := queryAllContainer()
		if len(results) == 0 {
			break
		}
		for _, result := range results {
			if containerListen[result.Key] == false  && result.Record.Used == "true"{
				stop[result.Key] = make(chan string,1)
				go startReciver(stop[result.Key],result.Key)
			}
			if containerListen[result.Key] == true  && result.Record.Used == "false" {
				stop[result.Key] <- "close"
			}
		}
		time.Sleep(time.Duration(2)*time.Second)
	}
}
