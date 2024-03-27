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

type Schneider struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewSchneider() *Schneider {
	return &Schneider{
		Brand:     "SCHNEIDER",
		PagesPath: "/var/www/simpleparse/pages/schneider/",
		CsvPath:   "files/csv/schneider.csv",
		Domain:    "schneider-russia.com",
		StartUrl:  "https://schneider-russia.com/silovoe-oborudovanie/predohraniteli/nojevye/predohranitel-1e-gf-200a",
	}
}

func (o *Schneider) Init() (bool, error) {
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

func (o *Schneider) OnLink(e *colly.HTMLElement) (bool, error) {
	link := e.Attr("href")

	// convert relative url to absolute
	url := e.Request.AbsoluteURL(link)

	if strings.Contains(url, "/kontrol-klimata/") {
		return true, nil
	}

	if strings.Contains(url, "/promyshlennaya-avtomatizaciya/") {
		return true, nil
	}

	if strings.Contains(url, "/silovoe-oborudovanie/") {
		return true, nil
	}

	return false, nil
}

func (o *Schneider) OnPage(e *colly.Response) (bool, error) {

	url := e.Request.URL.String()



	// if HTML contains div with class "product-page" then this is the page

// fmt.Println("CHECK IF PAGE", url)
	if !bytes.Contains(e.Body, []byte(`class="product-wrap"`)) {
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

func (o *Schneider) ParsePage(html *bytes.Buffer, source string) (*dto.Product, error) {

	// бренд, категория, модель, название, артикул, цена, источник, файл

	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	

	// parse sku from div.product-main-row ul.main-props-table li span.main-props-value 
	sku := ""
	doc.Find("div.product-main-row ul.main-props-table li span.main-props-value").Each(func(i int, s *goquery.Selection) {
			sku = strings.TrimSpace(s.Find("strong").Text())
	})

	// parse model from div.box-description text (whithout strong text) trim space
	model := ""

	// Find name from div.product-wrap h1 text trim space
	name := strings.TrimSpace(doc.Find("div.product-wrap h1").Text())

	// Find price from p.prod-price span.price-regular text trim "₽" trim space "₽"
	price := strings.TrimSpace(doc.Find("div.product-main-row p.prod-price span.price-regular").Text())
	price = strings.Replace(price, "₽", "", -1)
	price = strings.TrimSpace(price)



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
