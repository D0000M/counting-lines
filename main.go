package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
			continue
		}
		if isSelected(e.Name()) {
			filename := filepath.Join(dir, e.Name())
			count <- counting(filename)
		}
	}
}

func isSelected(name string) bool {
	suffix := getSuffix(*fileType)

	if filepath.Ext(name) != suffix {
		return false
	}
	name = strings.TrimSuffix(name, suffix)

	if *noTesting {
		cuts := strings.Split(name, "_")
		if cuts[len(cuts)-1] == "test" {
			return false
		}
	}
	return true
}

func getSuffix(fileType string) string {
	switch fileType {
	case "go":
		return ".go"
	case "c":
		return ".c"
	default:
		return ".go"
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

var fileType = flag.String("f", "go", "默认为go语言")
var noTesting = flag.Bool("t", false, "默认包含测试文件")

func main() {

	flag.Parse()

	numlines := start()
	fmt.Printf("一共%d行代码\n", numlines)
}
