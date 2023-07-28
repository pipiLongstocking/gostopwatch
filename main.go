package main

import "gostopwatch/examples"

func main() {
	err := examples.BasicTimer()
	if err != nil {
		panic(err)
	}
	err = examples.PauseTimer()
	if err != nil {
		panic(err)
	}
	err = examples.StopTimer()
	if err != nil {
		panic(err)
	}
}
