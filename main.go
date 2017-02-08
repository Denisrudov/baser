package main

import (
	"os"
	"path/filepath"
	"sync"
	"fmt"
	"baser/style"
	"flag"
)

const (
	IMAGE_SIZE int64 = 1024
)

var (
	imageSize int64
)

func init() {
	flag.Int64Var(&imageSize, "imsize", 0, "Size of the image in kbytes")
}

func main() {

	flag.Parse()

	if imageSize == 0 {
		imageSize = IMAGE_SIZE * 1024
	} else {
		imageSize = imageSize * 1024
	}

	args := flag.Args()

	files := []string{}

	rootPath := "."

	if len(args) == 1 {
		rootPath = os.Args[0]
	}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)

		if ext == ".css" || ext == ".scss" {
			files = append(files, path)

		}
		return err
	})

	if err != nil {
		fmt.Printf("Error found: %s\n", err.Error())
		os.Exit(0)
	}

	if len(files) < 1 {
		fmt.Println("No Files found.")
		os.Exit(0)
	}

	wg := sync.WaitGroup{}

	wg.Add(len(files))

	for _, aFile := range files {
		go func(aFile string) {
			if file, err := style.NewFile(aFile); err == nil {
				file.SetMaxImage(imageSize)
				err = file.Encode()
			}
			fmt.Printf("Processed: %s\n", aFile)

			wg.Done()
		}(aFile)
	}

	wg.Wait()

	fmt.Println("All files are processed!")

}
