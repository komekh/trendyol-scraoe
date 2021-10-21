package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Create another collector to scrape course details
	productCollector := c.Clone()

	productDetailCollector := c.Clone()

	c.OnHTML("div[id=navigation-wrapper]", func(e *colly.HTMLElement) {
		//fmt.Println(e)

		//e.ForEachWithBreak("a[href]", func(index int, elm *colly.HTMLElement) bool {
		//	fmt.Println(elm.ChildText(""))
		//	return true
		//})

		e.ForEach("a[href]", func(index int, e *colly.HTMLElement) {
			link := e.Attr("href")
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)

			productCollector.Visit(e.Request.AbsoluteURL(link))
		})
	})

	productCollector.OnHTML(".p-card-wrppr > .p-card-chldrn-cntnr", func(e *colly.HTMLElement) {
		//fmt.Println("product: ", e)
		e.ForEach("a[href]", func(index int, productElement *colly.HTMLElement) {
			//fmt.Println(productElement)
			link := productElement.Attr("href")
			//title := e.DOM.Find(".prdct-desc-cntnr-ttl").Text()
			//name := e.DOM.Find(".prdct-desc-cntnr-name").Text()

			//fmt.Printf("Link found: %q -> %s\n", productElement.Text, link)
			//fmt.Printf("Title: %s, name %s: link: %s \n", title, name, link)

			productDetailCollector.Visit(e.Request.AbsoluteURL(link))

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

		// ************************** Product Title, Name Section Begin ************************************
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

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.trendyol.com/")
}
