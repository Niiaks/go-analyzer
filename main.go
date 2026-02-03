package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const Size = 3 * 1 << 30 // 3 gigabytes

func main() {
	root := "C:/Users/junio"

	fmt.Printf("Scanning directory %s\n", root)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.Size() > Size {
			fmt.Printf("the file name is %s, located at %s with size %d\n", d.Name(), path, fileInfo.Size())
		}

		return nil

	})

	fmt.Println("files not large in current directory")

	if err != nil {
		log.Fatal(err)
	}
}

func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := sha256.New()

	//stream file contents into the hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}
