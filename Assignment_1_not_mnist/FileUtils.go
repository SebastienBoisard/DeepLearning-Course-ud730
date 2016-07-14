package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// downloadFile downloads a file if not present, and make sure it's the right size.
func downloadArchive(url string, fileName string, expectedBytes int64, force bool) {

	// Test if fileName exists or if the download must be done
	if _, err := os.Stat(fileName); os.IsNotExist(err) || force == true {
		// path/to/whatever does not exist

		response, err := http.Get(url + fileName)
		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(fileName, body, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	statInfo, err := os.Stat(fileName)

	if err != nil {
		log.Fatal(err)
	}

	fileSize := statInfo.Size()
	fmt.Println("file size=", statInfo.Size())

	if fileSize != expectedBytes {
		log.Fatal(fileName, " found but the size is different from what expected (expected ", expectedBytes, " != actual=", fileSize)
	}

	fmt.Println(fileName, " found and verified.")
}

// Extract a .tar.gz archive
func extractArchive(fileName string, force bool) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join("target", header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	const url = "http://commondatastorage.googleapis.com/books1000/"
	downloadArchive(url, "notMNIST_small.tar.gz", 8458043, false)
	err := extractArchive("notMNIST_small.tar.gz", false)
	if err != nil {
		log.Fatal(err)
	}
}
