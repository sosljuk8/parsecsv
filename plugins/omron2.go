package plugins

import (
	"fmt"
	"log"
	"parsecsv/dto"
	"parsecsv/orm"

	//"os"

	"github.com/xuri/excelize/v2"
)

type Omron2 struct {
	Brand     string
	PagesPath string
	CsvPath   string
	Domain    string
	StartUrl  string
}

func NewOmron2() *Omron2 {
	return &Omron2{
		Brand:     "OMRON",
		PagesPath: "files/xlsx/OMRON_RLP_2020_RUS_.xlsx",
		CsvPath:   "files/csv/omron2.csv",
		Domain:    "sensoren.ru",
		StartUrl:  "https://sensoren.ru/brands/omron/",
	}
}

func (o *Omron2) Init() (bool, error) {
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

func (o *Omron2) XlsxRead() {
	f, err := excelize.OpenFile(o.PagesPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("S1")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := [][]string{}

	for ir, row := range rows {
		if ir == 0 {
			continue
		}

		// if ir > 45 {
		// 	break
		// }
		//fmt.Println(ir, row[0])

		pr := &dto.Product{
			Brand:    o.Brand,
			Model:    row[1],
			Name:     row[6],
			SKU:      row[0],
			Price:    row[2],
			Currency: "RUB",
			Source:   "",
			File:     Hash(o.Brand) + ".html",
		}
		data = append(data, pr.String())

		// rm := map[int]string{}
		// 		 for ic, colCell := range row {

		// rm[ic] = string(colCell[ic])
		// 		// 	pr := &dto.Product{
		// 		// 		Brand:    o.Brand,
		// 		// 		Model:    string(colCell[1]),
		// 		// 		Name:     string(colCell[7]),
		// 		// 		SKU:      string(colCell[0]),
		// 		// 		Price:    string(colCell[2]),
		// 		// 		Currency: "RUB",
		// 		// 		Source:   "",
		// 		// 		File:     Hash(o.Brand) + ".html",
		// 		// 	}
		// 			//data = append(data, pr.String())
		// 		}
		// 		fmt.Println(rm)
	}
	//fmt.Println("ffffffffffffffff", len(data))
	err = orm.CreateAllFromSlice(o.CsvPath, data)

	//fmt.Fprintln(os.Stdout, []any{data}...)
}
