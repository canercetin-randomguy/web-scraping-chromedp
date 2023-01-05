package main

import (
	"canercetin/pkg/credentials"
	"canercetin/pkg/devprotocol"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// NavigateMainPage navigates the page that contains products and returns the html
func NavigateMainPage(navigationLink string, ctx *context.Context, buf *string, bodyPageXPath string) *string {
	err := chromedp.Run(*ctx,
		chromedp.Navigate(navigationLink),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.InnerHTML(bodyPageXPath, buf, chromedp.BySearch),
	)
	if err != nil {
		panic(err)
	}
	return buf
}

// NavigateProductPage navigates to the product page and returns the html
// productBodyXPath is the xpath of the product body
func NavigateProductPage(ctx *context.Context, link string, htmlBuffer *string, productBodyXPath string) *string {
	err := chromedp.Run(*ctx,
		chromedp.Navigate(link),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.OuterHTML("//*[(@id = \"aspnetForm\")]", htmlBuffer, chromedp.BySearch),
	)
	if err != nil {
		panic(err)
	}
	return nil
}

// NavigatePage navigates to the page and returns the html
//
// res == the html buffer
//
// csvFileName == Output's csv file name
//
// columnHeaderNames == the column names of the csv file
//
// baseLink == the base link of the page that will be paginated over, this is subjected to chanhge
//
// pageExtension == links extension, html?, aspx?, etc.
func NavigatePage(ctx *context.Context, res []byte, csvFileName string, baseLink string, pageExtension string, columnHeaderNames ...string) {
	// create [][]string to write to csv file
	var attemptCount = 1
	var data [][]string
	var productTitle string
	var productImage string
	var productSize string
	var productPrice string
	var productDescription string
	// start writing to csv file
	f, err := os.Create(csvFileName)
	f.Close()
	if err != nil {
		panic(err)
	}
	// create [][] string from the columnHeaderNames
	data = append(data, columnHeaderNames)
	if err != nil {
		return
	}
	var htmlBuffer string
	var soldOut = false
	// navigate to the page
	for pageCount := 1; pageCount <= 16; pageCount++ {
		navigationLink := baseLink + "p" + strconv.Itoa(pageCount) + ".aspx"
		fmt.Println(navigationLink)
		err = chromedp.Run(*ctx,
			chromedp.Navigate(navigationLink),
			chromedp.WaitReady("body", chromedp.ByQuery),
			chromedp.CaptureScreenshot(&res),
			// save
			chromedp.InnerHTML("//*[@id=\"bdyMaster\"]", &htmlBuffer, chromedp.BySearch),
		)
		if err != nil && strings.Contains(err.Error(), "net::ERR_ABORTED") == false {
			panic(err)
		} else if err != nil && strings.Contains(err.Error(), "net::ERR_ABORTED") == true {
			// wait for a random amount of time
			rand.Seed(time.Now().UnixNano())
			waitTime := rand.Intn(10)
			time.Sleep(time.Second * time.Duration(waitTime))
			attemptCount++
			if attemptCount > 15 {
				panic("too many attempts, yeter amk")
			}
			fmt.Println("connection aborted, trying again for the ", attemptCount, " time")
			NavigateMainPage(navigationLink, ctx, &htmlBuffer, "//*[@id=\"bdyMaster\"]")
		}
		attemptCount = 1
		// create io.Reader from htmlBuffer
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBuffer))
		if err != nil {
			panic(err)
		}
		// find all the links
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			class, exists := s.Attr("class")
			// nextProdName == product
			if exists == true && class == "nextProdName" {
				fmt.Println(s.Text())
				link, exists := s.Attr("href")
				if exists == true {
					// use the found link and navigate&scrape
					err := chromedp.Run(*ctx,
						chromedp.Navigate(link),
						chromedp.WaitReady("body", chromedp.ByQuery),
						chromedp.OuterHTML("//*[(@id = \"aspnetForm\")]", &htmlBuffer, chromedp.BySearch),
					)
					if err != nil && strings.Contains(err.Error(), "net::ERR_ABORTED") == false {
						panic(err)
					} else if err != nil && strings.Contains(err.Error(), "net::ERR_ABORTED") == true {
						// wait for a random amount of time
						rand.Seed(time.Now().UnixNano())
						waitTime := rand.Intn(10)
						time.Sleep(time.Second * time.Duration(waitTime))
						attemptCount++
						if attemptCount > 15 {
							panic("too many attempts, yeter amk")
						}
						fmt.Println("connection aborted, trying again for the ", attemptCount, " time")
						NavigateProductPage(ctx, link, &htmlBuffer, "//*[(@id = \"aspnetForm\")]")
					}
					doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBuffer))
					if err != nil {
						panic(err)
					}
					doc.Find("h1").Each(func(i int, s *goquery.Selection) {
						class, exists := s.Attr("class")
						if exists == true && class == "nextProdName" {
							id, exists := s.Attr("id")
							if exists == true && id == "ctl00_ctl00_ctl00_cphMain_cphMain_cphMain_pdtProduct_h1Name" {
								title, exists := s.Attr("itemprop")
								if exists == true && title == "name" {
									productTitle = s.Text()
								}
							}
						}
					})
					doc.Find("div").Each(func(i int, s *goquery.Selection) {
						property, exists := s.Attr("itemprop")
						if exists == true && property == "description" {
							productDescription = s.Text()
							fmt.Println(productDescription)
						}
					})
					doc.Find("div").Each(func(i int, s *goquery.Selection) {
						class, exists := s.Attr("class")
						if exists == true && class == "nextSoldOut" {
							productPrice = "Sold Out"
							soldOut = true
						}
					})
					doc.Find("img").Each(func(i int, s *goquery.Selection) {
						class, exists := s.Attr("class")
						if exists == true && class == "nextProdImage" {
							productImage, _ = s.Attr("src")
							fmt.Println(productImage)
						}
					})
					if soldOut == false {
						doc.Find("span").Each(func(i int, s *goquery.Selection) {
							property, exists := s.Attr("itemprop")
							if exists == true && property == "price" {
								fmt.Println(s.Text())
								productPrice = s.Text()
							}
						})
					}
					doc.Find("td").Each(func(i int, s *goquery.Selection) {
						class, exists := s.Attr("class")
						if exists == true && class == "nextCustomField4" {
							// delete all whitespaces in s.Text
							productCode := strings.ReplaceAll(s.Text(), " ", "")
							// delete all blank lines in s.Text
							productCode = strings.ReplaceAll(productCode, "\n", "")
							productSize = productCode
							fmt.Println(productSize)
						}
					})
					data = append(data, []string{productTitle, productDescription, productImage, productPrice, productSize})
				}
			}
			// put a random delay
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			soldOut = false
		})
		// write data to csv file
		f, err = os.OpenFile(csvFileName, os.O_APPEND|os.O_WRONLY, 0o644)

		w := csv.NewWriter(f)
		for _, record := range data {
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
		w.Flush()
		err = f.Close()
		if err != nil {
			return
		}
	}

}
func main() {
	fmt.Println("Welcome.")
	fmt.Println("Do we need to login? (y/n)")
	var login string
	fmt.Scanln(&login)
	// Buffer to hold the HTML
	var res []byte
	// make chromedp run headless
	taskCtx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	if login == "y" {
		cred := credentials.GetCredentials()
		devprotocol.LoginCTX(
			&taskCtx,
			res,
			cred.Username,
			cred.Password,
			// XPATH FIELDS
			cred.LoginUsernameField,
			cred.LoginPasswordField,
			cred.LoginButtonField,
			cred.LoginLink,
		)
	}
	cancel()
}
