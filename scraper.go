package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/fatih/color"
)

var re = regexp.MustCompile(`((https?:\/\/)(.*\.org))(\/[a-z]{1,}\/)`)

func scrape(url string) (media []string) {
	newurl := fetchRedirectedURL(url)
	board := extractBoard(newurl)
	response, err := http.Get(newurl)
	if err != nil {
		color.Red("Error loading thread: %s", err.Error())
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		color.Red("Error reading response: %s", err.Error())
		os.Exit(1)
	}

	var thread Posts
	err = json.Unmarshal(body, &thread)
	if err != nil {
		color.Red("Error reading JSON: %s", err.Error())
		color.Red("Thread has most likely 404'd!")
		os.Exit(1)
	}
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
	newurl := url + ".json"
	response, err := http.Get(newurl)
	if err != nil {
		color.Red("Error loading thread: %s", err.Error())
		os.Exit(1)
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
