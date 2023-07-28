package main

import "gostopwatch/examples"

func main() {
	err := examples.BasicTimer()
	if err != nil {
		panic(err)
	}
}
