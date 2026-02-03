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

func main() {

	// dir is the directory to start scanning from (flag "-dir", default: ".")
	// default scans in working directory
	dir := flag.String("dir", ".", "path to directory to begin scanning from")

	flag.Parse()

	fmt.Printf("Scanning files in dir %s\n", *dir)

	//a channel which handles all paths
	jobs := make(chan string, 200)

	var wg sync.WaitGroup

	workerCounter := 4

	for range workerCounter {
		wg.Add(1)

		wg.Go(func() {
			defer wg.Done()
			for path := range jobs {
				checkFile(path)
			}
		})
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
func checkFile(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.Size() > Size {
		fmt.Printf("Large file: %s | Size: %.2f GB\n",
			path,
			float64(info.Size())/(1<<30))
	}
}
