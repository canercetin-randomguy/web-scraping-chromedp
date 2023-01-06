package main

import (
	"canercetin/pkg/backend"
	"canercetin/pkg/credentials"
	"canercetin/pkg/devprotocol"
	"canercetin/pkg/logger"
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

func main() {
	go func() {
		// change backend.Port to your liking.
		err := backend.StartWebPageBackend(backend.Port)
		if err != nil {
			log.Println(err)
		}
	}()
	dbLogFile := logger.CreateNewFile("./logs/db")
	dbLogger, err := logger.NewLoggerWithFile(dbLogFile)
	if err != nil {
		log.Println(err)
	}
	// get a fresh database connection
	dbConn := sqlpkg.SqlConn{}
	err = dbConn.GetSQLConn("")
	if err != nil {
		dbLogger.Error(err.Error())
	}
	go func() {
		err = dbConn.CreateClientSchema()
		if err != nil {
			dbLogger.Error(err.Error())
		}
		err = dbConn.CreateClientTable()
		if err != nil {
			dbLogger.Error(err.Error())
		}
	}()
	fmt.Println("Welcome.")
	fmt.Println("Do we need to login? (y/n)")
	var login string
	fmt.Scanln(&login)
	// Buffer to hold the HTML
	var res string
	// make chromedp run headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancelAllocatedCTX := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAllocatedCTX()
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	// Ensure we started the browser
	if err := chromedp.Run(taskCtx); err != nil {
		log.Fatal(err)
	}
	if login == "y" {
		cred := credentials.GetCredentials()
		err := devprotocol.LoginCTX(
			&taskCtx,
			cred.Username,
			cred.Password,
			// XPATH FIELDS
			cred.LoginUsernameField,
			cred.LoginPasswordField,
			cred.LoginButtonField,
			cred.LoginLink,
		)
		if err != nil {
			fmt.Println("Uh oh.")
			panic(err)
		}
		fmt.Println("")
	} else {
		fmt.Println("Now lets see what do we have here...")
		// assign a dummy value to res
		// TODO: Yeet the colly lines to scraper.go
		// TODO: Let the user input the URL and find allowed domains from the URL
		// TODO: Create a simple UI
		res = " "
		_, err := devprotocol.NavigatePageReturnHTML("https://www.kktcarabam.com/kategori/ikinci-el-araclar-ticari-araclar-minibus-midibus",
			&taskCtx,
			&res,
			"body")
		if err != nil {
			fmt.Println("Uh oh.")
			panic(err)
		}
		c := colly.NewCollector(
			colly.AllowedDomains("www.bursadabugun.com", "bursadabugun.com"),
			colly.MaxDepth(2))
		var links []string
		var brokenLinks []string
		// Find and visit all links
		start := time.Now()
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			exists := false
			for _, link := range links {
				if link == e.Attr("href") {
					exists = true
				}
			}
			if exists == false {
				links = append(links, e.Attr("href"))
			}
			e.Request.Visit(e.Attr("href"))
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting -> ", r.URL)
		})

		c.OnResponse(func(r *colly.Response) {
			if r.StatusCode == 404 {
				brokenLinks = append(brokenLinks, r.Request.URL.String())
			}
		})

		c.Visit("https://www.bursadabugun.com/")
		fmt.Println("Total time taken: ", time.Since(start))
		fmt.Println("Total links found: ", len(links))
		// save broken links to a file
		file, err := os.Create("brokenLinks.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()
		for _, link := range brokenLinks {
			_, err := file.WriteString(link + "\n")
			if err != nil {
				log.Fatal("Cannot write to file", err)
			}
		}
		/*
			fmt.Println(*res)
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(*res))
			if err != nil {
				fmt.Println("Uh oh.")
				panic(err)
			}
			var availableLinks []string
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists == true {
					availableLinks = append(availableLinks, href)
				}
			})
		*/
	}
}
