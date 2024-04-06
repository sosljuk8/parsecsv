package plugins

import (
	"bytes"
	"fmt"
	"log"
	"parsecsv/dto"
	"parsecsv/orm"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Cisco struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewCisco() *Cisco {
	return &Cisco{
		Brand:     "CISCO",
		PagesPath: "/var/www/simpleparse/pages/cisco/",
		CsvPath:   "files/csv/cisco.csv",
		Domain:    "www.router-switch.com",
		StartUrl:  "https://www.router-switch.com/cisco-products-a-z.html",
	}
}

func (o *Cisco) Init() (bool, error) {
	// CreateDir(PagesPath)
	// err := CreateDir(o.PagesPath)
	// if err != nil {
	// 	log.Println(err)
	// 	return false, err
	// }

	err := CreateCsv(o.CsvPath)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (o *Cisco) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "/Price-cisco") {
		return true, nil
	}
	if strings.Contains(url, ".html") {
		return true, nil
	}
	return false, nil
}

func (o *Cisco) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()

	// if HTML contains div with class "product-page" then this is the page

	fmt.Println("CHECK IF PAGE", url)
	if !bytes.Contains(e.Body, []byte(`class="product-info-main-content pb-5"`)) {
		// just skip this url, no errors triggered
		return false, nil
	}

	product, err := o.ParsePage(bytes.NewBuffer(e.Body), url)
	if err != nil {
		return true, err
	}

	err = orm.WriteCsv(o.CsvPath, product)
	if err != nil {
		return true, err
	}

	// err = orm.SavePage(o.PagesPath+product.File, string(e.Body))
	// if err != nil {
	// 	return true, err
	// }

	return true, nil
}

func (o *Cisco) ParsePage(html *bytes.Buffer, source string) (*dto.Product, error) {

	// бренд, категория, модель, название, артикул, цена, источник, файл

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	// parse sku from table.product-data-table tbody  (meta[name="sku"] content)
	sku := strings.TrimSpace(doc.Find("table.product-data-table tbody meta[itemprop='sku']").AttrOr("content", ""))

	// parse model from table.product-data-table tbody tr td.h-model-pr
	model := strings.TrimSpace(doc.Find("table.product-data-table tbody td.h-model-pr").Text())

	// Find name from meta[name="keywords"] content
	name := strings.TrimSpace(doc.Find("table.product-data-table tbody div[itemprop='description']").Text())

	pPrice := doc.Find("table.product-data-table tbody td.p-listprice span.price").Text()
	sPrice := doc.Find("table.product-data-table tbody td.saleprice span.regular-price span.price").Text()

	// Find price div.price-novat text ; replece "Цена без НДС - " to ""; replace " руб." to "" and remove all spaces
	price := pPrice + " / " + sPrice

	currency := "USD"

	//fmt.Println(sku)

	return &dto.Product{
		Brand:    o.Brand,
		Model:    model,
		Name:     name,
		SKU:      sku,
		Price:    price,
		Currency: currency,
		Source:   source,
		File:     Hash(sku) + ".html",
	}, nil
}
