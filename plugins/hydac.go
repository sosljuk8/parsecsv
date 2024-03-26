package plugins

import (
	"bytes"
	"log"
	"parsecsv/dto"
	"parsecsv/orm"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Hydac struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewHydac() *Hydac {
	return &Hydac{
		Brand:     "HYDAC",
		PagesPath: "/var/www/simpleparse/pages/hydac/",
		CsvPath:   "files/csv/hydac.csv",
		Domain:    "www.hypneu.de",
		StartUrl:  "https://www.hypneu.de/shop/hydac-international-gmbh.html",
	}
}

func (o *Hydac) Init() (bool, error) {
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

func (o *Hydac) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "shop/hydac-international-gmbh") {
		return true, nil
	}
	return false, nil
}

func (o *Hydac) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()

	// if HTML contains div with class "product-page" then this is the page
	if !bytes.Contains(e.Body, []byte(`class="product-view"`)) {
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

func (o *Hydac) ParsePage(html *bytes.Buffer, source string) (*dto.Product, error) {

	// бренд, категория, модель, название, артикул, цена, источник, файл

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	occ := []string{
		"\n",
		"\t",
		"Description:",
		"Artikelnummer:",
		"Beschreibung:",
		"Datasheet",
		"PDF",
		"Product code:",
		"€",
	}

	// parse model from div.box-description text (whithout strong text) trim space
	model := strings.TrimSpace(doc.Find("div.box-description").Text())

	model = ReplaceString(model, occ)

	// // parse sku from div.product-shop first p text trim space
	sku := strings.TrimSpace(doc.Find("div.product-shop p").Eq(0).Text())
	sku = ReplaceString(sku, occ)



	// Find name from div.ajax_breadcrumbs h1 text trim space
	name := ""

	// Find price from div.price text trim ₽ trim space
	price := strings.TrimSpace(doc.Find("div.price-box span.price").Text())
	price = ReplaceString(price, occ)

	// Find currency from div.product-price-count meta itemprop="priceCurrency" (attribute "content" value)
	currency := "EUR"

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
