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
			if containerListen[result.ContainerID] == false  && result.Used == "true"{
				stop[result.ContainerID] = make(chan string,1)
				go startReciver(stop[result.ContainerID],result.ContainerID)
			}
			if containerListen[result.ContainerID] == true  && result.Used == "false" {
				stop[result.ContainerID] <- "close"
			}
		}
		time.Sleep(time.Duration(2)*time.Second)
	}
}
