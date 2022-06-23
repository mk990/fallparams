package headless

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chromedp/chromedp"
)

func Request(url string) string {
	var fullDOM string
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &fullDOM),
	)
	if err != nil {
		fmt.Println(err)
	}

	return fullDOM
}

func Screenshot(url string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		fmt.Println(err)
	}
	if err := ioutil.WriteFile(GenerateNameFromUrl(url)+".png", buf, 0644); err != nil {
		fmt.Println(err)
	}
}

func GenerateNameFromUrl(url string) string {
	url = strings.Replace(url, "://", "_", -1)
	url = strings.Replace(url, ".", "_", -1)
	url = strings.Replace(url, "/", "_", -1)
	url = strings.Replace(url, "?", "_", -1)
	url = strings.Replace(url, "=", "_", -1)
	url = strings.Replace(url, "&", "_", -1)
	url = strings.Replace(url, ":", "_", -1)
	url = strings.Replace(url, " ", "_", -1)
	return url
}
