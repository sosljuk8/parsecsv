package plugins

import (
	"bytes"
	"encoding/json"
	"log"
	"parsecsv/dto"
	"parsecsv/orm"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type SMC struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewSMC() *SMC {
	return &SMC{
		Brand:     "SMC",
		PagesPath: "/var/www/simpleparse/pages/smc/",
		CsvPath:   "files/csv/SMC.csv",
		Domain:    "industriation.ru",
		StartUrl:  "https://industriation.ru/smc/",
	}
}

func (o *SMC) Init() (bool, error) {
	CreateDir(o.PagesPath)
	err := CreateDir(o.PagesPath)
	if err != nil {
		log.Println(err)
		return false, err
	}

	err = CreateCsv(o.CsvPath)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (o *SMC) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "/smc/") {
		return true, nil
	}

	if strings.Contains(url, "-smc-") {
		return true, nil
	}

	return false, nil
}

func (o *SMC) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()

	// if HTML contains div with class "product-page" then this is the page
	if !bytes.Contains(e.Body, []byte(`<img src="https://industriation.ru/image/catalog/smc-products/logo_smc_corporation.png" title="SMC" alt="SMC">`)) {
		// just skip this url, no errors triggered
		return false, nil
	}

	product, err := o.ParsePage(bytes.NewBuffer(e.Body), url)
	if err != nil {
		return true, err
	}

	err = orm.WriteCsvP(o.CsvPath, product)
	if err != nil {
		return true, err
	}

	err = orm.SavePage(o.PagesPath+product.File, string(e.Body))
	if err != nil {
		return true, err
	}

	return true, nil
}

func (o *SMC) ParsePage(html *bytes.Buffer, source string) (*dto.PCard, error) {

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	card := dto.NewPCard()
	card.Brand = o.Brand
	card.Source = source

	// Brand
	// Category
	// Model
	// Name
	// SKU
	// Price
	// Currency
	// Source
	// Img
	// Properties
	// Description
	// File

	// Find img from div#images first img.image (attribute "src" value)
	card.Img = doc.Find("div#images img.image").First().AttrOr("src", "")

	// find description from div#description noindex first p text trim space
	//card.Description = strings.TrimSpace(doc.Find("div#description noindex p").First().Text())
	card.Description = ""

	// Find properties (map[string]string) from div#description div.arrt-line; key is the div.nma text, value is the div.txa text
	properties := make(map[string]string)
	doc.Find("div#description div.arrt-line").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("div.nma").Text())
		value := strings.TrimSpace(s.Find("div.txa").Text())
		properties[key] = value
	})

	// if properties is not empty then convert map[string]string
	if len(properties) > 0 {
		// properties = props to json
		propertiesJSON, err := json.Marshal(properties)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		card.Properties = string(propertiesJSON)
	}

	// Find sku from properties["Артикул"]
	card.SKU = strings.TrimSpace(properties["Артикул"])

	modelser := strings.TrimSpace(properties["Серия товара"])

	title := doc.Find("h1.heading-title").Text()

	sp := strings.Split(title, modelser)

	//model := modelser + sp.last item
	card.Model = modelser + sp[len(sp)-1]

	// Find name from div.ajax_breadcrumbs h1 text trim space
	card.Name = strings.TrimSpace(properties["Наименование"])

	// Find price from div.price text trim ₽ trim space
	price := strings.TrimSpace(doc.Find("div.product-data div.product-price-box div.price").Text())
	card.Price = strings.Trim(price, "₽")

	// Find currency from div.product-price-count meta itemprop="priceCurrency" (attribute "content" value)
	card.Currency = "RUB"
	card.File = Hash(card.SKU) + ".html"

	//cat, desc := ParseCat(card.SKU)

	card.Category = ""
	card.Description = ""

	return card, nil
}
