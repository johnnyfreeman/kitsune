package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func tailFile(wg *sync.WaitGroup, filename string) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("[ERROR] Cannot open file: %s (%v)", filename, err)
		return
	}
	defer file.Close()

	// Seek to end of file
	file.Seek(0, io.SeekEnd)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(500 * time.Millisecond) // Wait for new data
			continue
		}
		fmt.Printf("[%s] %s", filename, line) // Prefix output with filename
	}
}

func main() {
	flag.Parse()
	files := flag.Args()

	if len(files) == 0 {
		log.Fatal("Usage: go run multitail.go <file1> <file2> ...")
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go tailFile(&wg, file)
	}

	wg.Wait()
}
