package main

import (
	"gpomo/gostopwatch"
	"log"
	"time"
)

func TestStopwatch() {
	p, err := gostopwatch.NewGStopwatch(3 * time.Second)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	endTimer := false
	for {
		select {
		case x := <-p.Tick:
			log.Printf("channel op: %+v", x)
		case <-p.Done:
			log.Println("timer ended")
			endTimer = true
		}
		if endTimer {
			break
		}
	}
	log.Println("end test")
}

func TestStopwatchPause() {
	log.Println("Start pause test")
	p, err := gostopwatch.NewGStopwatch(5 * time.Second)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(1010 * time.Millisecond)
	err = p.Pause()
	if err != nil {
		panic(err)
	}
	log.Printf("time left %v", p.GetTimeLeft())
	time.Sleep(2 * time.Second)
	err = p.Resume()
	if err != nil {
		panic(err)
	}
	endTimer := false
	for {
		select {
		case x := <-p.Tick:
			log.Printf("channel op: %+v", x)
		case <-p.Done:
			log.Println("timer ended")
			endTimer = true
		}
		if endTimer {
			break
		}
	}
	log.Println("end test")
}

func main() {
	TestStopwatch()
	TestStopwatchPause()

}
