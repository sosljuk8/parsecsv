package orm

import (
	"fmt"
	"os"
	"strings"
)

func SavePage(filename string, html string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		return err
	}

	return nil

}

func Unmapped(d string) ([]string, error){
	f, err := os.Open(d)
    if err != nil {
        fmt.Println(err)
        return  nil, err
    }
    files, err := f.Readdir(0)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

	unm := []string{}

    for _, v := range files {

		if v.IsDir(){
			continue
		}
		if strings.Contains(v.Name(), "MAP"){
			continue
		}

		unm = append(unm, v.Name())
    }

return unm, nil
}