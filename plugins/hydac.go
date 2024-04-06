package plugins

import (
	"bytes"
	"encoding/json"
	"log"
	"parsecsv/dto"
	"parsecsv/orm"
	"regexp"
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
		PagesPath: "/var/www/brands/pages/hydac/",
		CsvPath:   "/var/www/brands/csv/hydac.csv",
		Domain:    "www.hydac.com",
		StartUrl:  "https://www.hydac.com/shop/en/hps-2400-1000496612",
	}
}

func (o *Hydac) Init() (bool, error) {
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

func (o *Hydac) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "shop/en") {
		return true, nil
	}
	return false, nil
}

func (o *Hydac) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()

	// if HTML contains div with class "product-page" then this is the page
	if !regexp.MustCompile(`\/shop\/en\/\d+$`).MatchString(url) {
		// just skip this url, no errors triggered
		return false, nil
	}

	product, err := o.ParsePage(bytes.NewBuffer(e.Body), url)
	if err != nil {
		return true, err
	}

	err = orm.WritePCsv(o.CsvPath, product)
	if err != nil {
		return true, err
	}

	err = orm.SavePage(o.PagesPath+product.File, string(e.Body))
	if err != nil {
		return true, err
	}

	return true, nil
}

func (o *Hydac) ParsePage(html *bytes.Buffer, source string) (*dto.PCard, error) {

	// бренд, категория, модель, название, артикул, цена, источник, файл

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	// parse model from h1.page-title span.base attribute "data-product-name"
	model := doc.Find("h1.page-title span.base").Text()

	// parse sku
	sku := doc.Find(".sku .value").Text()

	// parse category
	// Find category from div.breadcrumbs ul (li a text +|+ li a text+|+li a text)
	category := ""
	doc.Find(".breadcrumbs ul li a").Each(func(i int, s *goquery.Selection) {

		category += "|" + strings.TrimSpace(s.Text())
	})
	category = strings.Replace(category, "|Home|", "", -1)

	// parse img from div.product first img.fotorama__img (attribute "data-src" value)
	img := doc.Find("div.product img.lazyload").First().AttrOr("data-src", "")

	// Find name from div.ajax_breadcrumbs h1 text trim space
	name := ""

	// Find price from div.price text trim ₽ trim space
	price := ""

	// parse properties
	props := make(map[string]string)
	doc.Find(".additional-attributes-wrapper table tbody tr").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find(".label").Text())
		value := s.Find(".data").Text()
		props[key] = value
	})

	properties := ""

	// if properties is not empty then convert map[string]string
	if len(props) > 0 {
		// properties = props to json
		propertiesJSON, err := json.Marshal(props)
		if err != nil {
			return nil, err
		}
		properties = string(propertiesJSON)
	}

	description := ""

	// Find currency from div.product-price-count meta itemprop="priceCurrency" (attribute "content" value)
	currency := "EUR"

	//fmt.Println(sku)

	return &dto.PCard{
		Brand:       o.Brand,
		Category:    category,
		Model:       model,
		Name:        name,
		SKU:         sku,
		Price:       price,
		Currency:    currency,
		Source:      source,
		Img:         img,
		Properties:  properties,
		Description: description,
		File:        Hash(sku) + ".html",
	}, nil
}
