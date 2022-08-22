package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	c := make(chan string)
	var totalRead uint64
	var totalWrite uint64
	var wg sync.WaitGroup
	maxMessage := uint64(50000000)

	startTime := time.Now().Unix()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := uint64(0); i < maxMessage; i++ {
			c <- "foo"
			totalWrite++
		}
		fmt.Println(fmt.Sprintf("write Done, totalWrite %d", totalWrite))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			<-c
			totalRead++
			if totalRead == maxMessage {
				return
			}
		}
	}()

	wg.Wait()
	endTime := time.Now().Unix()
	fmt.Println(fmt.Sprintf("done total Read %d, startTime %d, endTime %d", totalRead, startTime, endTime))

}
