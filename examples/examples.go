package examples

import (
	"gostopwatch/gostopwatch"
	"log"
	"time"
)

// BasicTimer is an implementation of basic timer functionality.
func BasicTimer() error {
	log.Println("BasicTimer begin")
	sw, err := gostopwatch.NewGStopwatch(5 * time.Second)
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
			log.Println("BasicTimer end")
			return nil
		}
	}
}

// PauseTimer shows how to pause/resume the timer.
func PauseTimer() error {
	log.Println("PauseTimer begin")
	sw, err := gostopwatch.NewGStopwatch(5 * time.Second)
	if err != nil {
		return err
	}
	tick, done, err := sw.Start()
	if err != nil {
		return nil
	}
	go func() {
		time.Sleep(2 * time.Second)
		errP := sw.Pause()
		if errP != nil {
			panic(err)
		}
		log.Println("Paused()")
		time.Sleep(2 * time.Second)
		errP = sw.Resume()
		if errP != nil {
			panic(err)
		}
		log.Println("Resumed()")

	}()
	for {
		select {
		case t := <-tick:
			log.Printf("tick: %v", t)
		case <-done:
			log.Println("PauseTimer end")
			return nil
		}
	}
}

// StopTimer shows how to stop the timer.
func StopTimer() error {
	log.Println("StopTimer begin")
	sw, err := gostopwatch.NewGStopwatch(5 * time.Second)
	if err != nil {
		return err
	}
	tick, done, err := sw.Start()
	if err != nil {
		return nil
	}
	go func() {
		time.Sleep(2 * time.Second)
		errP := sw.Stop()
		if errP != nil {
			panic(err)
		}
		log.Println("Stopped()")
	}()
	for {
		select {
		case t := <-tick:
			log.Printf("tick: %v", t)
		case <-done:
			log.Println("StopTimer end")
			return nil
		}
	}

}
