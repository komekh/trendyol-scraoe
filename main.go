package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"scrap/db"
	"scrap/models"
	"strings"
)

func init() {
	fmt.Println("init fn called")
	db.Setup()
}

func main() {

	categories := make([]models.Category, 0)

	// Instantiate default collector
	c := colly.NewCollector()

	//storage := &postgres.Storage{
	//	URI:          "host=localhost port=5432 user=postgres password=admin database=scrap sslmode=disable",
	//	VisitedTable: "colly_visited",
	//	CookiesTable: "colly_cookies",
	//}
	//
	//if err := c.SetStorage(storage); err != nil {
	//	panic(err)
	//}

	// Collector to scrape products
	productCollector := c.Clone()

	// Collector for product details
	productDetailCollector := c.Clone()

	// Collector for Categories
	//categoryCollector := c.Clone()

	idCounter := 1
	currParentId := idCounter

	// Iterate all Categories
	c.OnHTML(".main-nav > .tab-link", func(e *colly.HTMLElement) {
		//fmt.Println("e: ", e)
		//link := e.Attr("href")
		//fmt.Println("Link: ", link)
		var currSubParentId int
		e.ForEach("a[href]", func(index int, catHeaderElm *colly.HTMLElement) {
			link := catHeaderElm.Attr("href")

			// Category header types
			// 1. category-header => main category. L1
			// 2. sub-category header -> L2
			// 3. L3 category does not exist class
			categoryType := catHeaderElm.Attr("class")

			category := models.Category{
				Id:       idCounter,
				Name:     catHeaderElm.Text,
				Link:     link,
				ParentId: 0,
				Level:    0,
			}
			if categoryType == "category-header" {
				//category.Level = 0
				currParentId = idCounter

			} else if categoryType == "sub-category-header" {
				category.Level = 1
				category.ParentId = currParentId
				currSubParentId = idCounter
			} else {
				category.Level = 2
				category.ParentId = currSubParentId
			}

			sqlQuery := `INSERT INTO categories (id, name, link, parentId, level )
				 VALUES ($1, $2, $3, $4, $5);`
			_, err := db.Sqlx.Db.Exec(sqlQuery, category.Id, category.Name, category.Link, category.ParentId, category.Level)

			if err != nil {
				fmt.Println("ERROR: ", err)
				//panic(err)
			}

			categories = append(categories, category)
			//fmt.Printf("i: %d, Category: %s, Link:%s\n", index, category.Name, category.Link)
			idCounter++
		})
		fmt.Println("***********************")
	})

	// Iterate all Categories
	//c.OnHTML("div[id=navigation-wrapper]", func(e *colly.HTMLElement) {

	//// Iterate subcategories
	//e.ForEach("a[href]", func(index int, e *colly.HTMLElement) {
	//	link := e.Attr("href")
	//	fmt.Printf("Link found: %d: %q -> %s\n", index, e.Text, link)
	//
	//	productCollector.Visit(e.Request.AbsoluteURL(link))
	//})
	//})

	productCollector.OnHTML(".p-card-wrppr > .p-card-chldrn-cntnr", func(e *colly.HTMLElement) {
		//fmt.Println("product: ", e)
		e.ForEach("a[href]", func(index int, productElement *colly.HTMLElement) {
			//fmt.Println(productElement)
			link := productElement.Attr("href")

			//title := e.DOM.Find(".prdct-desc-cntnr-ttl").Text()
			//name := e.DOM.Find(".prdct-desc-cntnr-name").Text()

			fmt.Printf("Link found: %q -> %s\n", productElement.Text, link)
			//fmt.Printf("Title: %s, name %s: link: %s \n", title, name, link)

			//productDetailCollector.Visit(e.Request.AbsoluteURL(link))

		})
		/* WORKING.... DONT DELETE
		fmt.Println("---------- PRODUCT BEGIN ------------------")
		title := e.DOM.Find(".prdct-desc-cntnr-ttl").Text()
		name := e.DOM.Find(".prdct-desc-cntnr-name").Text()
		link := e.Attr("href")
		fmt.Printf("Title: %s, name %s: link: %s \n", title, name, link)
		fmt.Println("---------- PRODUCT END --------------------")
		fmt.Println()
		*/
	})

	productDetailCollector.OnHTML(".product-container", func(e *colly.HTMLElement) {

		// ************************** Product Title, Name Section Begin ***********************************
		productHeader := e.ChildText("h1.pr-new-br")
		productName := e.ChildText("h1.pr-new-br > span")
		productTitle := strings.ReplaceAll(productHeader, productName, "")
		fmt.Println("Product Title: ", productTitle, "Product Name: ", productName)
		// ************************** Product Title, Name Section End ************************************

		// ************************** Price Section **************************************************
		// 1. original price
		// 2. discount in cart
		// 3. discount without cart

		// 1. Original Price
		priceOriginal := e.ChildText("div.product-price-container > div.pr-bx-w > div.pr-bx-nm > span.prc-slg")

		// 2. Discount in Cart
		priceDiscount := ""
		discountDesc := ""
		if len(priceOriginal) == 0 {
			priceOriginal = e.ChildText("span.prc-slg.prc-slg-w-dsc")
			priceDiscount = e.ChildText("span.prc-dsc")
			discountDesc = e.ChildText("div.pr-bx-pr-dsc > .pr-bx-pr-dsc")
		}

		// 3. Discount without Cart. Original price does not change
		discountedStamp := ""
		priceDiscount = e.ChildText("div.pr-bx-w > div.pr-bx-nm with-org-prc > span.prc-slg")
		discountedStamp = e.ChildText("div.discounted-stamp > span.discounted-stamp-text")

		fmt.Println("priceOriginal: ", priceOriginal)
		fmt.Println("PriceDiscount: ", priceDiscount)
		fmt.Println("DiscountDesc: ", discountDesc)
		fmt.Println("DiscountedStamp: ", discountedStamp)
		// ************************ Price Section End **************************************************

		// ************************ Merchant Section **************************************************
		//Nullable
		merchant := e.ChildText("a.merchant-text")
		fmt.Println("Merchant: ", merchant)
		// ************************ Merchant Section End ***********************************************

		// ************************ Color Section ***********************************************
		e.ForEach(".sp-itm", func(i int, sizeElm *colly.HTMLElement) {
			fmt.Println("size: ", sizeElm.Text)
		})
		// ************************ Color Section End ***********************************************
	})

	// TODO:
	// 1. Colors
	// 2. Image URL
	// 3. Options

	//c.OnHTML(".main-nav", func(e *colly.HTMLElement) {
	//	fmt.Println(e)
	//})

	// On every element which has href attribute call callback
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	link := e.Attr("href")
	//	// Print link
	//	fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	//	// Visit link found on page
	//	// Only those links are visited which are in AllowedDomains
	//	c.Visit(e.Request.AbsoluteURL(link))
	//})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	//c.OnScraped(func(response *colly.Response) {
	//	for _, category := range categories {
	//		fmt.Printf("ID: %d, name: %s, parentID: %d, level: %d\n", category.Id, category.Name, category.ParentId, category.Level)
	//	}
	//})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.trendyol.com/")
}
