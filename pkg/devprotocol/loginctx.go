package devprotocol

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
	"math/rand"
	"strings"
	"time"
)

// LoginCTX logs into the website and returns the context for using it in the future actions
//
// ctx is the context to be used
//
// buf is the buffer to be used
//
// usernameFormXPath is the xpath of the username form
//
// passwordFormXPath is the xpath of the password form
//
// loginButtonXPath is the xpath of the login button
//
// siteLink is the link of the website
//
// username and password is self explanatory
func LoginCTX(ctx *context.Context, username string, password string, usernameFormXPath string, passwordFormXPath string, loginButtonXPath string, siteLink string) error {
	var attemptCount = 1
	err := chromedp.Run(*ctx,
		chromedp.Navigate(siteLink),
		chromedp.WaitReady("body"),
		chromedp.Click(usernameFormXPath, chromedp.BySearch),
		chromedp.Sleep(time.Second*1),
		chromedp.SetValue(usernameFormXPath, username, chromedp.BySearch),
		chromedp.Click(passwordFormXPath, chromedp.BySearch),
		chromedp.Sleep(time.Second*1),
		chromedp.SetValue(passwordFormXPath, password, chromedp.BySearch),
		chromedp.Click(loginButtonXPath, chromedp.BySearch),
		chromedp.Sleep(time.Second*1),
	)

	// ChromeDP may fail to log in with ERR_ABORTED, so we try again for 15 times. If it fails, we panic. Paniiiiic.
	if err != nil && strings.Contains(err.Error(), "net::ERR_ABORTED") == false {
		panic(err)
	} else if err != nil && strings.Contains(err.Error(), "net::ERR_ABORTED") == true {
		rand.Seed(time.Now().UnixNano())
		// wait between 0-3 seconds.
		waitTime := rand.Intn(3)
		time.Sleep(time.Second * time.Duration(waitTime))
		attemptCount++
		if attemptCount > 15 {
			panic("too many attempts, yeter uLa")
		}
		fmt.Println("connection aborted, trying again for the ", attemptCount, " time")
		LoginCTX(ctx, username, password, usernameFormXPath, passwordFormXPath, loginButtonXPath, siteLink)
	}
	return err
}
