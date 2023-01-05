package devprotocol

import (
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
)

func NavigatePageReturnHTML(navigationLink string, ctx *context.Context, buf *string, bodyPageXPath string) (*string, error) {
	err := chromedp.Run(*ctx,
		chromedp.Navigate(navigationLink),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.InnerHTML(bodyPageXPath, buf, chromedp.BySearch),
	)
	return buf, err
}
func NavigatePage(navigationLink string, ctx *context.Context) {
	err := chromedp.Run(*ctx,
		chromedp.Navigate(navigationLink),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
	if err != nil {
		panic(err)
	}
}
