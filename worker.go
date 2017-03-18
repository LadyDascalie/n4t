package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

var fileRegExp = regexp.MustCompile(`[0-9]`)

func download(media string, wg *sync.WaitGroup, bar *pb.ProgressBar) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()
	defer wg.Done()
	if !silent {
		defer bar.Increment()
	}

	resp, err := http.Get(media)
	if err != nil {
		fails.Get++
	}

	filename := fileRegExp.FindAllString(media, -1)  // find file id
	filename = append(filename, filepath.Ext(media)) // extract file extension

	fn := strings.Join(filename, "") // Add it all together

	file, err := os.Create(fn)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fails.Copy++
	}
}
