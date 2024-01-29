package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func counting(filename string) int {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "counting(): %v\n", err)
		return 0
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var count int
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "counting(): %v\n", err)
		return 0
	}

	return count
}

func dirents(dir string) []fs.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dirents: %v\n", err)
		return nil
	}
	return entries
}

func walkDir(dir string, n *sync.WaitGroup, count chan<- int) {
	defer n.Done()
	entrys := dirents(dir)
	for _, e := range entrys {
		if e.IsDir() {
			n.Add(1)
			subDir := filepath.Join(dir, e.Name())
			go walkDir(subDir, n, count)
		} else {
			if filepath.Ext(e.Name()) != ".go" {
				continue
			}
			filename := filepath.Join(dir, e.Name())
			count <- counting(filename)
		}
	}
}

func start() int {
	count := make(chan int)
	var n sync.WaitGroup

	n.Add(1)
	go walkDir(`./`, &n, count)

	go func() {
		n.Wait()
		defer close(count)
	}()

	var numlines int
	for c := range count {
		numlines += c
	}

	return numlines
}

func main() {
	fmt.Println(start())
}
