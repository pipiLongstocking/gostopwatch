package main

import (
	"gostopwatch/gostopwatch"
	"log"
	"time"
)

func timer() error {
	sw, err := gostopwatch.NewGStopwatch(10 * time.Second)
	if err != nil {
		return err
	}
	tick, done, err := sw.Start()
	if err != nil {
		return nil
	}
	for {
		select {
		case t := <-tick:
			log.Printf("tick: %v", t)
		case <-done:
			log.Println("timer ended")
			return nil
		}
	}
}

func main() {
	err := timer()
	if err != nil {
		panic(err)
	}
}
