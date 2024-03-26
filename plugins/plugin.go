package plugins

import (
	"encoding/hex"
	"hash/fnv"
	"log"
	"os"
	"strings"
)

func Hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateCsv(path string) error {

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	return nil
}

func ReplaceString(s string, occ []string) string {

	for _, oldString := range occ {
		s = strings.Replace(s, oldString, "", -1)
	}

	return strings.TrimSpace(strings.Replace(s, "\n", "", -1))
}
