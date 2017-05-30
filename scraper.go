package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
)

// This regex breaks down the 4chan urls into components
// https://regex101.com/r/NQIQ4H/1/
var re = regexp.MustCompile(`((https?:\/\/)(.*\.org))(\/[a-z]{1,}\/)`)
var threadUrlWithoutUserInput = regexp.MustCompile(`(.*)(\/thread\/)([0-9]{1,})`)

func scrape(url string) (media []string) {
	// get the redirected url
	postRedirectionURL := fetchRedirectedURL(url)
	// get the board letter
	board := extractBoard(postRedirectionURL)

	// Fetch the post
	response, err := http.Get(postRedirectionURL)
	if err != nil {
		color.Red("Error loading thread: %s", err.Error())
		os.Exit(1)
	}

	// Read it
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		color.Red("Error reading response: %s", err.Error())
		os.Exit(1)
	}

	// Unmarshall it
	var thread Posts
	err = json.Unmarshal(body, &thread)
	if err != nil {
		spew.Dump(err)
		color.Red("Error reading JSON: %s", err.Error())
		color.Red("Thread has most likely 404'd!")
		os.Exit(1)
	}

	// Create the media urls
	for _, post := range thread.Posts {
		if post.Tim != 0 {
			media = append(media, fmt.Sprintf("https://%s%s%d%s", cdnStem, board, post.Tim, post.Ext))
		}
	}
	return media
}

func extractBoard(url string) string {
	matches := re.FindAllStringSubmatch(url, -1)
	return matches[0][4]
}

func fetchRedirectedURL(url string) string {
	matches := threadUrlWithoutUserInput.FindAllStringSubmatch(url, -1)
	urlToThreadAsJSON := matches[0][0] + ".json"
	response, err := http.Get(urlToThreadAsJSON)
	if err != nil {
		color.Red("Error loading thread: %s", err.Error())
	}
	return response.Request.URL.String()
}

// Posts is the full thread object
type Posts struct {
	Posts []Post `json:"posts,omitempty"`
}

// Post is a single post within the thread
type Post struct {
	Tim int    `json:"tim,omitempty"`
	Ext string `json:"ext,omitempty"`
}
