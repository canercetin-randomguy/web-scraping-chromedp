package links

import (
	"canercetin/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
	"log"
	"net/url"
	"strings"
)

// IsUrl  Thanks a lot to https://stackoverflow.com/a/55551215/17996217 for the code.
//
// Checks if url is valid or not.
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// FindLinks is used for collecting whole links from the website.
//
// Example usage: FindLinks("www.google.com", 2, "canercetin")
//
// canercetin parameter means that client "canercetin" initiated the request.
//
// At the end of the process, returns links and broken links in order.
//
// Returns (nil,nil) on error.
//
// Check terminal if something went wrong while creating the logger.
func FindLinks(siteLink string, maxDepth int, username string, linkLimit int) string {
	// make a seperate links and brokenLinks slice, self explanatory.
	var ScrapedLinks = make(map[string]LinkStruct)
	var temporaryLink = LinkStruct{
		// Just dummy values, these will be changed and appended to ScrapedLinks.LinkStorage
		Link:     "example.com",
		IsBroken: false,
	}
	// Create a new folder for the client if it does not exist.
	err := logger.CreateNewFolder(username)
	if err != nil {
		log.Println(err)
		return ""
	}
	// Create a new logger to store the errors
	// fileNumber means, if we have a file called collector_canercetin_20220101_1 or collector_canercetin_20220101_0
	// fileNumber is 1 or 0. This increases when client requests a lot of stuff.
	collectorLogFile, fileNumber := logger.CreateNewFileCollector(fmt.Sprintf("./logs/%s/", username), username)
	collectorLogger, err := logger.NewLoggerWithFile(collectorLogFile)
	if err != nil {
		log.Println(err)
		return ""
	}

	errorLogFile, _ := logger.CreateNewFileError(fmt.Sprintf("./logs/%s/", username), username)
	errorLogger, err := logger.NewLoggerWithFile(errorLogFile)

	if err != nil {
		log.Println(err)
		return ""
	}
	// find the absolute path of the link, such as convert http://example.com to example.com, then store it in a seperate string.
	absoluteSiteLink := ConvertToAbsolutePath(siteLink)
	absoluteAbsoluteSitelink := strings.ReplaceAll(absoluteSiteLink, "www.", "")
	// So we will have a www.example.com and an example.com absolute domains to allow.
	c := colly.NewCollector(
		// allow specific domains, so we don't fucking yeet to cloudflare all of a sudden
		colly.AllowedDomains(absoluteSiteLink, absoluteAbsoluteSitelink),
		// let the customer decide how deep they want to go
		colly.MaxDepth(maxDepth))
	// let us wander in all a[href] tags, I mean they are considered links, aren't they?
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Check if valid first.
		if IsUrl(e.Attr("href")) {
			// Set a flag to false
			exists := false
			// Wander around in already collected links
			for it := range ScrapedLinks {
				// If the link is already in the slice, set the flag to true
				if ScrapedLinks[it].Link == e.Attr("href") {
					exists = true
				}
			}
			// If the flag is still false, append the link to the slice
			if exists == false {
				temporaryLink.Link = e.Attr("href")
				temporaryLink.IsBroken = false
				ScrapedLinks[e.Attr("href")] = temporaryLink
			}
		} else {
			// it may be something like /computers. So we need to add the domain to it.
		}
		err := c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		if err != nil {
			if strings.Contains(err.Error(), "already visited") {
				// do nothing
			} else {
				errorLogger.Errorw("Something went wrong while visiting the link.", zap.Error(err),
					"client", username,
					"fileNumber", fileNumber,
					"link", e.Attr("href"))
			}
		}
	})
	// Log when we request something, I mean, come on we have 20 GB space in cloud.
	c.OnRequest(func(r *colly.Request) {
		collectorLogger.Infow(fmt.Sprintf("Visiting %s", r.URL.String()),
			"url", r.URL.String(),
			"client", username,
			"fileNumber", fileNumber)
	})
	// Append the broken links to the brokenLinks slice
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 404 {
			if entry, ok := ScrapedLinks[r.Request.URL.String()]; ok {
				entry.IsBroken = true
				ScrapedLinks[r.Request.URL.String()] = entry
			}
		}
	})
	// Start wandering
	err = c.Visit(siteLink)
	if err != nil {
		errorLogger.Infow("Error while visiting the site",
			"error", err,
			"website", siteLink,
			"client", username)
	}
	// marshal the ScrapedLinks
	fmt.Println(ScrapedLinks)
	linkResponse, err := json.Marshal(ScrapedLinks)
	if err != nil {
		errorLogger.Errorw("Error while marshalling the links",
			"error", err,
			"client", username)
	}
	return string(linkResponse)
}

func ConvertToAbsolutePath(siteLink string) string {
	absoluteSiteLink := strings.ReplaceAll(siteLink, "http://", "")
	absoluteSiteLink = strings.ReplaceAll(absoluteSiteLink, "https://", "")
	// delete everything after /, such as example.com/abc/def/ghi to example.com
	absoluteSiteLink = strings.Split(absoluteSiteLink, "/")[0]
	return absoluteSiteLink
}
