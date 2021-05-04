package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Printf("Process ID: %d", os.Getpid())

	signalchan := make(chan os.Signal, 1)
	signal.Notify(signalchan)
	block := make(chan bool, 1)

	go func() {
		log.Println("Waiting for signal")
		log.Println("Press ctrl-c to stop")
		s := <-signalchan
		log.Printf("Received signal: %+v", s)
		block <- true
	}()

	<-block
}
