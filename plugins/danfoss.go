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

type Danfoss struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewDanfoss() *Danfoss {
	return &Danfoss{
		Brand:     "DANFOSS",
		PagesPath: "/var/www/simpleparse/pages/danfoss/",
		CsvPath:   "files/csv/danfoss.csv",
		Domain:    "dan-service.ru",
		StartUrl:  "https://dan-service.ru/catalog/",
	}
}

func (o *Danfoss) Init() (bool, error) {
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

func (o *Danfoss) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "/catalog/") {
		return true, nil
	}
	return false, nil
}

func (o *Danfoss) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()



	// if HTML contains div with class "product-page" then this is the page

fmt.Println("CHECK IF PAGE", url)
	if !bytes.Contains(e.Body, []byte(`class="catalog_detail detail element_1"`)) {
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

func (o *Danfoss) ParsePage(html *bytes.Buffer, source string) (*dto.Product, error) {

	// бренд, категория, модель, название, артикул, цена, источник, файл

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	

	// parse sku from div.right_info div.article  span.value text trim space
	sku := strings.TrimSpace(doc.Find("div.right_info div.article span.value").Text())

	// parse model from div.box-description text (whithout strong text) trim space
	model := ""

	// Find name from div.top_info div.preview_text text trim space
	name := strings.TrimSpace(doc.Find("div.top_info div.preview_text").Text())

	// Find price from div.prices_block div.price(currency= attr "data-currency" ; value= attr "data-value")  
	price := strings.TrimSpace(doc.Find("div.prices_block div.price").AttrOr("data-value", ""))
	currency := strings.TrimSpace(doc.Find("div.prices_block div.price").AttrOr("data-currency", ""))

	

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
