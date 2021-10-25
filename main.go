package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"scrap/db"
	"scrap/scrapers"
)

func init() {
	fmt.Println("init fn called")
	db.Setup()
}

func main() {

	fmt.Println("Hello World!")

	scraper := scrapers.Scraper{
		Collector: colly.NewCollector(),
	}

	//scraper.CategoryScrapper()
	//scraper.ProductScraper()
	scraper.BrandScraper()
	fmt.Println("Read From DB")
}
