package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const Size = 3 * 1 << 30 // 3 gigabytes

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go-analyzer [-d directory] [-s size]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	// dir is the directory to start scanning from (flag "-dir", default: ".")
	// default scans in working directory
	dir := flag.String("d", ".", "path to directory to begin scanning from")

	// minSize is the minimum file size threshold (in gigabytes) used to decide whether a file
	// should be reported as "large". Set with -s; files with size > minSize are reported.
	minSize := flag.Int("s", 1, "minimum file size in GB")
	minBytes := int64(*minSize) << 30

	flag.Usage = usage
	flag.Parse()

	fmt.Printf("Scanning files in dir %s\n", *dir)

	//a channel which handles all paths
	jobs := make(chan string, 200)

	var wg sync.WaitGroup

	workerCounter := 4

	for range workerCounter {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for path := range jobs {
				checkFile(path, minBytes)
			}
		}()
	}

	err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			jobs <- path
		}

		return nil

	})
	close(jobs)
	wg.Wait()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scan complete.")
}

// checkFile prints the file and it's size
func checkFile(path string, size int64) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.Size() > size {
		fmt.Printf("Large file: %s | Size: %.2f GB\n",
			path,
			float64(info.Size())/(1<<30))
	}
}
