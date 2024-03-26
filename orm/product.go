package orm

import (
	"encoding/csv"
	"log"
	"os"
	"parsecsv/dto"
)

func WriteCsv(path string, product *dto.Product) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(product.String())

	return nil
}

func CreateAllFromSlice(path string, products [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
    err = w.WriteAll(products) // calls Flush internally

    if err != nil {
        log.Fatal(err)
    }

	return nil
}

