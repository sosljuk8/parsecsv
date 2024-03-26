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

type Ifm struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewIfm() *Ifm {
	return &Ifm{
		Brand:     "IFM Electronic",
		PagesPath: "/var/www/simpleparse/pages/ifm_electronic/",
		CsvPath:   "files/csv/ifm.csv",
		Domain:    "sensoren.ru",
		StartUrl:  "https://sensoren.ru/brands/ifm_electronic/",
	}
}

func (o *Ifm) Init() (bool, error) {
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

func (o *Ifm) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "ifm_electronic") {
		return true, nil
	}

	return false, nil
}

func (o *Ifm) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()

	// if HTML contains div with class "product-page" then this is the page
	if !bytes.Contains(e.Body, []byte(`class="product-page"`)) {
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

func (o *Ifm) ParsePage(html *bytes.Buffer, source string) (*dto.Product, error) {

	// бренд, категория, модель, название, артикул, цена, источник

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	// Find product skuStr from div.product-data first div.product-data-elem first trim space
	skuStr := doc.Find("div.product-data div.product-data-elem").First()
	// find sku from skuStr second span
	sku := skuStr.Find("span").Eq(1).Text()
	sku = strings.TrimSpace(sku)
	model := sku

	// Find name from div.ajax_breadcrumbs h1 text trim space
	name := strings.TrimSpace(doc.Find("div.ajax_breadcrumbs h1").Text())

	// Find price from div.product-price-count-actual (attribute "content" value)
	price, _ := doc.Find("div.product-price-count-actual").Attr("content")

	// Find currency from div.product-price-count meta itemprop="priceCurrency" (attribute "content" value)
	currency, _ := doc.Find("div.product-price-count meta[itemprop=priceCurrency]").Attr("content")

	//sku = strings.TrimSpace(sku)
	fmt.Println(sku, model, name, price, currency, source)

	//fmt.Println(sku)

	return &dto.Product{
		Brand:    o.Brand,
		Model:    model,
		Name:     name,
		SKU:      sku,
		Price:    price,
		Currency: currency,
		Source:   source,
		File:     Hash(model) + ".html",
	}, nil
}
