package main

import (
	"canercetin/pkg/credentials"
	"canercetin/pkg/devprotocol"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
	"log"
	"strings"
)

func main() {
	fmt.Println("Welcome.")
	fmt.Println("Do we need to login? (y/n)")
	var login string
	fmt.Scanln(&login)
	// Buffer to hold the HTML
	var res *string
	// make chromedp run headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
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
		res = devprotocol.NavigatePageReturnHTML("https://www.kktcarabam.com/kategori/ikinci-el-araclar-ticari-araclar-minibus-midibus",
			&taskCtx,
			res,
			"body")
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(*res))
		if err != nil {
			fmt.Println("Uh oh.")
			panic(err)
		}
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists == true {
				fmt.Println(href)
			}
		})
	}
	cancel()
}
