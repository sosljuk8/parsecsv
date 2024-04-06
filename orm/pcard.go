package orm

import (
	"encoding/csv"
	"log"
	"os"
	"parsecsv/dto"
)

func WritePCsv(path string, product *dto.PCard) error {
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
