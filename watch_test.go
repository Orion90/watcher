package watcher

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createFile() {
	err := os.MkdirAll("import", 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	createFile()
	d := "import"
	defer os.RemoveAll("import") // clean up
	c, cancel := New("import")
	content := []byte("temporary file's content")
	tmpfn := filepath.Join(d, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, content, 0777); err != nil {
		log.Fatal(err)
	}
	go func() {
		select {
		case <-time.After(1 * time.Second):
			cancel()
			break
		}
	}()
	for f := range c {
		log.Println("TEST: ", f)
	}
}
