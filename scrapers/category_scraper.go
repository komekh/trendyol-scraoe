package scrapers

import (
	"fmt"
	"github.com/gocolly/colly"
	"scrap/helper"
	"scrap/models"
)

type Scraper struct {
	Collector *colly.Collector
}

func (scraper *Scraper) CategoryScrapper() {
	categories := make([]models.Category, 0)

	idCounter := 1
	currParentId := idCounter

	// Iterate all Categories and save to DB
	scraper.Collector.OnHTML(".main-nav > .tab-link", func(e *colly.HTMLElement) {
		var currSubParentId int
		e.ForEach("a[href]", func(index int, catHeaderElm *colly.HTMLElement) {
			link := catHeaderElm.Attr("href")

			// Category header types
			// 1. category-header => main category. L1
			// 2. sub-category-header -> L2
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

			//sqlQuery := `INSERT INTO categories (id, name, link, parentId, level )
			//	 VALUES ($1, $2, $3, $4, $5);`
			//_, err := db.Sqlx.Db.Exec(sqlQuery, category.Id, category.Name, category.Link, category.ParentId, category.Level)
			//
			//if err != nil {
			//	fmt.Println("ERROR: ", err)
			//	//panic(err)
			//}

			categories = append(categories, category)
			//fmt.Printf("i: %d, Category: %s, Link:%s\n", index, category.Name, category.Link)
			idCounter++
		})
		fmt.Println("***********************")
	})

	scraper.Collector.OnScraped(func(response *colly.Response) {
		fmt.Println("Category scrape completed")
		for _, category := range categories {
			fmt.Printf("ID: %d, name: %s, parentID: %d, level: %d\n", category.Id, category.Name, category.ParentId, category.Level)
		}
	})

	// Start scraping on https://www.trendyol.com/
	scraper.Collector.Visit("https://www.trendyol.com/")
}

func (scraper *Scraper) ProductScraper() {
	//FIXME: solve product model
	product := models.Product{}
	//products := make([]models.Product, 0)

	productCollector := scraper.Collector.Clone()
	productDetailCollector := scraper.Collector.Clone()

	// Iterate all Categories
	scraper.Collector.OnHTML(".main-nav > .tab-link", func(e *colly.HTMLElement) {
		// Iterate subcategories
		e.ForEach("a[href]", func(index int, e *colly.HTMLElement) {
			link := e.Attr("href")
			//fmt.Printf("Link found: %d: %q -> %s\n", index, e.Text, link)

			//Visit products
			productCollector.Visit(e.Request.AbsoluteURL(link))
		})
	})

	// Collects Product from List
	productCollector.OnHTML(".p-card-wrppr", func(e *colly.HTMLElement) {
		fmt.Println("----Product Collector BEGIN-----")

		product.Id = e.Attr("data-id")

		e.ForEach("a[href]", func(index int, productElement *colly.HTMLElement) {
			product.Link = productElement.Attr("href")

			product.Title = e.DOM.Find(".prdct-desc-cntnr-ttl").Text()
			product.Name = e.DOM.Find(".prdct-desc-cntnr-name").Text()
			// FIXME: Get variant count
			product.ColorVariantCount = e.DOM.Find("span.color-variant-count").Text()

			//fmt.Println("COLOR COUNT: ", product.ColorVariantCount)

			// Visit product detail
			productDetailCollector.Visit(e.Request.AbsoluteURL(product.Link))

		})

		fmt.Println("----Product Collector END-----")
		fmt.Println()
	})

	productDetailCollector.OnHTML(".product-container", func(e *colly.HTMLElement) {

		fmt.Println("****Product DETAIL Collector BEGIN****")

		// ************************** Product Title, Name Section Begin ***********************************
		//productHeader := e.ChildText("h1.pr-new-br")
		//productName := e.ChildText("h1.pr-new-br > span")
		//productTitle := strings.ReplaceAll(productHeader, productName, "")
		//fmt.Println("Product Title: ", productTitle, "Product Name: ", productName)
		// ************************** Product Title, Name Section End ************************************

		// ************************** Price Section **************************************************
		// 1. original price
		// 2. discount in cart
		// 3. discount without cart

		// 1. Original Price
		product.PriceOrg = e.ChildText("div.product-price-container > div.pr-bx-w > div.pr-bx-nm > span.prc-slg")

		// 2. Discount in Cart
		product.PriceDisc = ""
		product.PriceDiscDesc = ""
		if len(product.PriceOrg) == 0 {
			product.PriceOrg = e.ChildText("span.prc-slg.prc-slg-w-dsc")
			product.PriceDisc = e.ChildText("span.prc-dsc")
			product.PriceDiscDesc = e.ChildText("div.pr-bx-pr-dsc > .pr-bx-pr-dsc")
		}

		// 3. Discount without Cart. Original price does not change
		product.PriceDiscStamp = ""
		product.PriceDisc = e.ChildText("div.pr-bx-w > div.pr-bx-nm with-org-prc > span.prc-slg")
		product.PriceDiscStamp = e.ChildText("div.discounted-stamp > span.discounted-stamp-text")

		// ************************ Price Section End **************************************************

		// ************************ Merchant Section **************************************************
		//Nullable
		//merchant := e.ChildText("a.merchant-text")
		//fmt.Println("Merchant: ", merchant)
		// ************************ Merchant Section End ***********************************************

		// ************************ Color Section ***********************************************
		e.ForEach(".sp-itm", func(i int, sizeElm *colly.HTMLElement) {
			product.Sizes = append(product.Sizes, sizeElm.Text)
		})
		// ************************ Color Section End ***********************************************

		// Get stamp image and push to images
		stampImg := e.ChildAttr("img.product-stamp", "src")
		if len(stampImg) != 0 {
			img := models.Image{
				BaseImage: false,
				IsStamp:   true,
				Link:      stampImg,
			}
			product.Images = append(product.Images, img)
		}

		// Get all images belong to product
		productImages := e.ChildAttrs("img", "src")

		// The last image is product's base image with size 1200x1800
		baseProductImg := productImages[len(productImages)-1]
		baseImg := models.Image{
			BaseImage: true,
			IsStamp:   false,
			Link:      baseProductImg,
		}
		product.Images = append(product.Images, baseImg)

		// remove stamp image and base product image
		helper.Remove(productImages, stampImg)
		helper.Remove(productImages, baseProductImg)

		for _, img := range productImages {
			pImg := models.Image{
				BaseImage: false,
				IsStamp:   false,
				Link:      img,
			}
			product.Images = append(product.Images, pImg)
		}

		// TODO:
		// 1. Colors
		// 2. Image URL
		// 3. Options

		//fmt.Println("PRODUCT:  ", product.Images)

		//fmt.Println("Product")
		//fmt.Println("Title: ", product.Title, "Name: ", product.Name, "LINK: ", product.Link, "ID: ", product.Id)
		//for i, image := range product.Images {
		//	fmt.Println("i: ", i, "image link: ", image.Link, "isBase: ", image.BaseImage, "isStamp: ", image.IsStamp)
		//}
		//
		//for i, size := range product.Sizes {
		//	fmt.Println("i: ", i, "Size: ", size)
		//}

		fmt.Println("****Product DETAIL Collector END****")
		//fmt.Println("****************************")
	})

	//productDetailCollector.OnHTML(".gallery-container", func(element *colly.HTMLElement) {
	//	el := element.Attr("img[src]")
	//	fmt.Println("EL: ", el)
	//})

	scraper.Collector.Visit("https://www.trendyol.com/")
}
