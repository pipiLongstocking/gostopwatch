package main

import (
	"gpomo/gostopwatch"
	"log"
	"time"
)

func main() {
	p, err := gostopwatch.NewGStopwatch(3 * time.Second)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	for {
		select {
		case x := <-p.Tick:
			if x == 0 {
				break
			}
		case <-p.Done:
			log.Printf("timer ended")
		}
	}
}
