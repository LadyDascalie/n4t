package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
	"github.com/ladydascalie/sortdir/sortdir"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	boardStem = "4chan.org"
	cdnStem   = "i.4cdn.org"
	dlFolder  = "n4t_downloads"
)

var (
	fails           Failures // Global fail count
	subFolderName   string   // flag
	threadURL       string   // flag
	silent          bool     // flag
	threadsOverride int

	// worker config
	wg sync.WaitGroup
	// Threads determines how many concurrent downloads are allowed
	Threads   = 10
	semaphore = make(chan struct{}, Threads)
)

func main() {
	flag.StringVar(&subFolderName, "f", "", "Choose a subfolder name:\n\t n4t -f folder_name")
	flag.StringVar(&threadURL, "u", "", "Choose a subfolder name:\n\t n4t -u thread_url")
	flag.IntVar(&threadsOverride, "t", 0, "Choose how concurrent downloads to run (max 12):\n\t n4t -t 5")
	flag.BoolVar(&silent, "s", false, "Choose silent output:\n\t n4t -s")
	flag.Parse()

	if threadsOverride > 0 && threadsOverride <= 12 {
		Threads = threadsOverride
		color.Green("Starting with %d concurrent downloads...", Threads)
	}

	var media []string

	// Get url then scrape it
	switch threadURL == "" {
	case true:
		url := getUserInput()
		media = scrape(url)
	case false:
		media = scrape(threadURL)
	}

	// Start the progress bar
	var bar *pb.ProgressBar
	count := len(media)
	if !silent {
		bar = pb.StartNew(count)
	}

	// Set the download location
	location := setDownloadLocation()

	for _, m := range media {
		wg.Add(1)
		go download(m, &wg, bar) // worker.go
	}

	wg.Wait()
	close(semaphore)

	if !silent {
		bar.FinishPrint(color.GreenString("%s", "Download completed!"))
	}

	// Prepare to sort by extension
	files := sortdir.Ls(location, true)
	sortdir.SortByTypes(files)

	if fails.Get > 0 || fails.Copy > 0 {
		color.Red("%s", fails.String())
	}
}

var failed int

func getUserInput() string {
	if failed > 1 {
		color.Yellow("%s", "Are you having problems pasting in the url?")
		color.Green("%s", "\tPlease try [ctrl+shift+v] (windows/linux), [cmd+v] (macOs) or right click, then select paste")
		os.Exit(1)
	}

	var url string
	notice := color.GreenString("%s", "Paste thread URL, then press 'Enter':")

	fmt.Println(notice)
	_, err := fmt.Scanln(&url)
	if err != nil {
		failed++
		color.Red("%s", "Invalid URL provided. Please confirm the URL then try again.")
		return getUserInput()
	}
	if govalidator.IsURL(url) {
		return url
	}
	failed++
	color.Red("%s", "Invalid URL provided. Please confirm the URL then try again.")
	return getUserInput()
}

// setDownloadLocation sets the download folder in the user's home folder
func setDownloadLocation() (downloadLocation string) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	downloadLocation = filepath.Join(usr.HomeDir, dlFolder)

	if subFolderName != "" {
		downloadLocation = filepath.Join(downloadLocation, subFolderName)
	}

	os.MkdirAll(downloadLocation, 0755)
	os.Chdir(downloadLocation)

	return downloadLocation
}
