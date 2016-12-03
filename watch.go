package watcher

import (
	"context"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

//New yields a new goroutine which watches path for file system changes and sends them on a channel
func New(path string) (<-chan string, context.CancelFunc) {
	log.Printf("Watching %v for created files", path)
	ctx, cancel := context.WithCancel(context.Background())
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan string)
	go loop(ctx, watcher, path, c)
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	return c, cancel
}

//loop runs in its own goroutine and sends events on create events
func loop(ctx context.Context, w *fsnotify.Watcher, path string, c chan string) {
	name := ""
	for {
		timer := time.NewTimer(10 * time.Second)
		select {
		case event := <-w.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				name = event.Name
			}
		case err := <-w.Errors:
			log.Println("error:", err)
		case <-ctx.Done():
			close(c)
			return
		case <-timer.C:
			if name != "" {
				c <- name
				name = ""
			}
		}
		timer.Stop()
	}
}
