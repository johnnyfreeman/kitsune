package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type LogFile struct {
	path   string
	mutex  sync.Mutex
	file   *os.File
	offset int64
}

func (lf *LogFile) Open() error {
	file, err := os.Open(lf.path)
	if err != nil {
		return err
	}
	lf.file = file
	return nil
}

func (lf *LogFile) ReadNewLines() ([]string, error) {
	lf.mutex.Lock()
	defer lf.mutex.Unlock()

	if lf.file == nil {
		return nil, fmt.Errorf("file not opened")
	}

	info, err := lf.file.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() < lf.offset {
		// File was truncated, reset offset
		lf.offset = 0
	}

	lf.file.Seek(lf.offset, 0)

	var lines []string
	scanner := bufio.NewScanner(lf.file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	lf.offset, _ = lf.file.Seek(0, os.SEEK_CUR)

	return lines, scanner.Err()
}

func watchFile(lf *LogFile, events chan<- string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(lf.path)
	if err != nil {
		log.Fatalf("Failed to watch %s: %v", lf.path, err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				lines, err := lf.ReadNewLines()
				if err == nil {
					for _, line := range lines {
						events <- fmt.Sprintf("[%s] %s", lf.path, line)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <logfile1> <logfile2> ...", os.Args[0])
	}

	logFiles := make([]*LogFile, 0)
	for _, path := range os.Args[1:] {
		lf := &LogFile{path: path}
		if err := lf.Open(); err != nil {
			log.Printf("Failed to open %s: %v", path, err)
		} else {
			logFiles = append(logFiles, lf)
		}
	}

	events := make(chan string, 100)
	var wg sync.WaitGroup

	for _, lf := range logFiles {
		wg.Add(1)
		go func(lf *LogFile) {
			defer wg.Done()
			watchFile(lf, events)
		}(lf)
	}

	go func() {
		for event := range events {
			// TODO: buffer events in the future
			fmt.Println(event)
		}
	}()

	wg.Wait()
}
