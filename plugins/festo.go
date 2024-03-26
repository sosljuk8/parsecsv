package plugins

import (
	"bytes"
	"log"
	"parsecsv/dto"
	"parsecsv/orm"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Festo struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewFesto() *Festo {
	return &Festo{
		Brand:     "FESTO",
		PagesPath: "/var/www/simpleparse/pages/festo/",
		CsvPath:   "files/csv/festo.csv",
		Domain:    "industriation.ru",
		StartUrl:  "https://industriation.ru/festo/",
	}
}

func (o *Festo) Init() (bool, error) {
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

func (o *Festo) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "/festo/") {
		return true, nil
	}
	if regexp.MustCompile(`industriation\.ru/\d+/`).MatchString(url) {
		return true, nil
	}
	return false, nil
}

func (o *Festo) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()

	// if HTML contains div with class "product-page" then this is the page
	if !bytes.Contains(e.Body, []byte(`<img src="https://industriation.ru/image/catalog/logo/brand1x.png" title="Festo" alt="Festo">`)) {
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

func (o *Festo) ParsePage(html *bytes.Buffer, source string) (*dto.Product, error) {

	// бренд, категория, модель, название, артикул, цена, источник, файл

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	// Find properties (map[string]string) from div#description div.arrt-line; key is the div.nma text, value is the div.txa text
	properties := make(map[string]string)
	doc.Find("div#description div.arrt-line").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("div.nma").Text())
		value := strings.TrimSpace(s.Find("div.txa").Text())
		properties[key] = value
	})

	// Find sku from properties["Артикул"]
	sku := strings.TrimSpace(properties["Артикул"])

	modelser := strings.TrimSpace(properties["Серия товара"])

	title := doc.Find("h1.heading-title").Text()

	sp := strings.Split(title, modelser)

	//model := modelser + sp.last item
	model := modelser + sp[len(sp)-1]

	// Find name from div.ajax_breadcrumbs h1 text trim space
	name := strings.TrimSpace(properties["Наименование"])

	// Find price from div.price text trim ₽ trim space
	price := strings.TrimSpace(doc.Find("div.product-data div.product-price-box div.price").Text())
	price = strings.Trim(price, "₽")

	// Find currency from div.product-price-count meta itemprop="priceCurrency" (attribute "content" value)
	currency := "RUB"

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
