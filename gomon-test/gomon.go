package main

import (
	"log"
	"os"
	"strings"

	filepath "path"

	"github.com/rjeczalik/notify"
)

func main() {
	var err error
	path := "."

	if len(os.Args) != 1 {
		path = os.Args[1]
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !strings.HasPrefix(path, "/") {
		path = filepath.Join(wd, path)
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Printf("Path does not exist: %s", path)
		panic(err)
	}

	log.Printf("Watching %s", path)

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening for events within a directory tree rooted
	// at current working directory. Dispatch remove events to c.
	if err := notify.Watch(path, c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	ei := <-c

	log.Println("Got event:", ei)
}
