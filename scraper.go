package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func scrape(url string) (media []string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	doc.Find("div.fileText a").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")

		if strings.Contains(url, boardStem) {
			media = append(media, fmt.Sprintf("https:%s", url))
		}
	})
	return media
}
