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

type Siemens struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewSiemens() *Siemens {
	return &Siemens{
		Brand:     "SIEMENS",
		PagesPath: "/var/www/simpleparse/pages/siemens/",
		CsvPath:   "files/csv/siemens.csv",
		Domain:    "simecs.ru",
		StartUrl:  "https://simecs.ru/catalog/siemens/avtomatizatsia_Siemens/sitop/6ep13333ba108ac0/",
	}
}

func (o *Siemens) Init() (bool, error) {
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

func (o *Siemens) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "catalog/siemens/") {
		return true, nil
	}
	return false, nil
}

func (o *Siemens) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()



	// if HTML contains div with class "product-page" then this is the page

fmt.Println("CHECK IF PAGE", url)
	if !bytes.Contains(e.Body, []byte(`class="row product-first-row"`)) {
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

func (o *Siemens) ParsePage(html *bytes.Buffer, source string) (*dto.Product, error) {

	// бренд, категория, модель, название, артикул, цена, источник, файл

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	

	// parse sku from div.product-wrap h1.bx-title text
	sku := strings.TrimSpace(doc.Find("div.product-wrap h1.bx-title").Text())


	// parse model from div.box-description text (whithout strong text) trim space
	model := ""

	// Find name from meta[name="keywords"] content
	name := strings.TrimSpace(doc.Find("meta[name='description']").AttrOr("content", ""))
	

	// Find price div.price-novat text ; replece "Цена без НДС - " to ""; replace " руб." to "" and remove all spaces
	price := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(doc.Find("div.price-novat").Text(), "Цена без НДС - ", ""), " руб.", ""))


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
