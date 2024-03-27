package main

import (
	"fmt"
	"os"
	"parsecsv/plugins"

	"github.com/gocolly/colly/v2"
)

// "bytes"
// "encoding/csv"
// "fmt"
// "log"
// "os"
// "parsecsv/dto"
// "strings"

// "github.com/PuerkitoBio/goquery"
// "github.com/gocolly/colly/v2"

// func init() {

// }

func main() {
	plugin := plugins.NewSchneider()

	_, err := plugin.Init()
	if err != nil {
		os.Exit(1)
	}

	Crawling(plugin)
	//XlsxRead(plugin)
}

func Maping() {
	// pages, err := orm.Unmapped(plugin.PagesPath)
	// if err != nil {
	// 	log.Println(err)
	// }

	// if len(pages) == 0 {
	// 	log.Println("Mapping Done")
	// 	os.Exit(0)
	// }

	// for _, v := range pages {
	// 	plugin.ParsePage(v, )
	// }
}

func XlsxRead(p *plugins.Lenze) {
	p.XlsxRead()
}

func Crawling(p *plugins.Schneider) {

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains(p.Domain),
		//colly.Async(),
	)

	//c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		// convert relative url to absolute
		url := e.Request.AbsoluteURL(link)

		visit, err := p.OnLink(e)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !visit {
			return
		}

		// Visit link found on page on a new thread
		c.Visit(url)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// After making a request print "Visited ..."
	c.OnResponse(func(r *colly.Response) {

		page, err := p.OnPage(r)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !page {
			return
		}

		fmt.Println("THIS IS THE PAGE!!!!", r.Request.URL)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(p.StartUrl)
	// Wait until threads are finished
	//c.Wait()
}
