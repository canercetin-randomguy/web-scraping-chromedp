package links

import (
	"canercetin/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
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
// At the end of the process, returns links as a json and fileNumber which represents clients file number.
//
// Returns ("",0) on error.
//
// Check terminal if something went wrong while creating the logger.
func FindLinks(siteLink string, maxDepth int, username string, linkLimit int) (string, int) {
	// make a seperate links and brokenLinks slice, self explanatory.
	var ScrapedLinks = make(map[string]LinkStruct)
	var temporaryLink = LinkStruct{
		// Just dummy values, these will be changed and appended to ScrapedLinks.LinkStorage
		Link:          "example.com",
		IsBroken:      false,
		OutsideDomain: false,
	}
	// Create a new folder for the client for logging if it does not exist.
	err := logger.CreateNewFolder(fmt.Sprintf("./logs/%s", username))
	if err != nil {
		log.Println(err)
		return "", 0
	}
	// Create a new logger to store the errors
	// fileNumber means, if we have a file called collector_canercetin_20220101_1 or collector_canercetin_20220101_0
	// fileNumber is 1 or 0. This increases when client requests a lot of stuff.

	errorLogFile, fileNumber := logger.CreateNewFileError(fmt.Sprintf("./logs/%s/", username), username)
	errorLogger, err := logger.NewLoggerWithFile(errorLogFile)
	if err != nil {
		log.Println(err)
		return "", 0
	}

	// find the absolute path of the link, such as convert www.example.com to example.com, then store it in a seperate string.
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
		if IsUrl(e.Request.AbsoluteURL(e.Attr("href"))) {
			// Set a flag to false
			exists := false
			// Wander around in already collected links
			for it := range ScrapedLinks {
				// If the link is already in the slice, set the flag to true
				if ScrapedLinks[it].Link == e.Request.AbsoluteURL(e.Attr("href")) {
					exists = true
				}
			}
			// If the flag is still false, append the link to the slice
			if exists == false {
				temporaryLink.Link = e.Request.AbsoluteURL(e.Attr("href"))
				temporaryLink.IsBroken = false
				temporaryLink.StatusCode = 404
				ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))] = temporaryLink
			}
		}
		err = c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		if err != nil {
			if strings.Contains(err.Error(), "already visited") {
				// do nothing
			} else if strings.Contains(err.Error(), "Not Found") {
				if entry, ok := ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))]; ok {
					entry.IsBroken = true
					entry.OutsideDomain = false
					entry.StatusCode = e.Response.StatusCode
					ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))] = entry
				}
			} else if strings.Contains(err.Error(), "Forbidden domain") {
				if entry, ok := ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))]; ok {
					if IsUrl(e.Request.AbsoluteURL(e.Attr("href"))) {
						req, _ := http.NewRequest(http.MethodGet, e.Request.AbsoluteURL(e.Attr("href")), nil)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							log.Println(err)
						}
						respData, _ := io.ReadAll(resp.Body)
						if string(respData) == "" {
							entry.IsBroken = true
							entry.OutsideDomain = true
							entry.StatusCode = 404
							ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))] = entry
						} else if resp.StatusCode != 200 {
							entry.IsBroken = true
							entry.OutsideDomain = true
							entry.StatusCode = resp.StatusCode
							ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))] = entry
						} else {
							entry.IsBroken = false
							entry.OutsideDomain = true
							entry.StatusCode = resp.StatusCode
							ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))] = entry
						}
					} else {
						entry.IsBroken = true
						entry.OutsideDomain = true
						entry.StatusCode = 404
						ScrapedLinks[e.Request.AbsoluteURL(e.Attr("href"))] = entry
					}
				}
			} else {
				errorLogger.Errorw("Something went wrong while visiting the link.", zap.Error(err),
					"client", username,
					"fileNumber", fileNumber,
					"link", e.Attr("href"))
			}
		}

	})
	// Append the broken links to the brokenLinks slice
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			if entry, ok := ScrapedLinks[r.Request.URL.String()]; ok {
				entry.IsBroken = true
				entry.OutsideDomain = false
				entry.StatusCode = r.StatusCode
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
	linkResponse, err := json.Marshal(ScrapedLinks)
	// convert the json to csv
	if err != nil {
		errorLogger.Errorw("Error while marshalling the links",
			"error", err,
			"client", username)
	}
	return string(linkResponse), fileNumber
}

func ConvertToAbsolutePath(siteLink string) string {
	absoluteSiteLink := strings.ReplaceAll(siteLink, "http://", "")
	absoluteSiteLink = strings.ReplaceAll(absoluteSiteLink, "https://", "")
	// delete everything after /, such as example.com/abc/def/ghi to example.com
	absoluteSiteLink = strings.Split(absoluteSiteLink, "/")[0]
	return absoluteSiteLink
}
