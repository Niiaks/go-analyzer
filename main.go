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

type Job struct {
	path string
	d    fs.DirEntry
}

func main() {

	// dir is the directory to start scanning from (flag "-dir", default: ".")
	// default scans in working directory
	dir := flag.String("d", ".", "path to directory to begin scanning from")

	// minSize is the minimum file size threshold (in gigabytes) used to decide whether a file
	// should be reported as "large". Set with -s; files with size > minSize are reported.
	minSize := flag.Float64("s", 1, "minimum file size in GB")

	flag.Usage = usage
	flag.Parse()

	minBytes := int64(*minSize) << 30
	fmt.Fprintf(os.Stdout, "Scanning files in dir %s\n", *dir)

	//a channel which handles all paths

	jobs := make(chan Job, 200)

	var wg sync.WaitGroup

	workerCounter := 4

	for range workerCounter {

		wg.Go(func() {
			for job := range jobs {
				checkFile(job.path, job.d, minBytes)
			}
		})
	}

	err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		job := Job{
			path: path,
			d:    d,
		}

		if !d.IsDir() {
			jobs <- job
		}

		return nil

	})
	close(jobs)
	wg.Wait()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stdout, "Scan complete.")
}

// checkFile prints the file and it's size
func checkFile(path string, d fs.DirEntry, size int64) {
	info, err := d.Info()
	if err != nil {
		return
	}

	if info.Size() > size {
		fmt.Fprintf(os.Stdout, "Large file: %s | Size: %.2f GB\n",
			path,
			float64(info.Size())/(1<<30))
	}
}
