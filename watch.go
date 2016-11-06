package watcher

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
)

//New yields a new goroutine which watches path for file system changes and sends them on a channel
func New(path string) (<-chan string, context.CancelFunc) {
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
	for {
		select {
		case event := <-w.Events:
			log.Println("event:", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				c <- event.Name
				log.Println("New file:", event.Name)
			}
		case err := <-w.Errors:
			log.Println("error:", err)
		case <-ctx.Done():
			close(c)
			return
		}
	}
}
