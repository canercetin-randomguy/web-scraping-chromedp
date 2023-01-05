package links

import (
	"github.com/PuerkitoBio/goquery"
)

func FindLinks(doc *goquery.Document, linkStorage map[string]bool) {
	// create a stopwatch to measure the time it takes to find the links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// check if the link is already in the map
			_, mapExists := linkStorage[href]
			if !mapExists {
				linkStorage[href] = true
			}
		}
	})
}
