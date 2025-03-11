package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/huh"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"table-separator/pkg/ods"
	"time"
)

var wg sync.WaitGroup

func main() {
	t := time.Now()

	fileName := flag.String("f", "", "(ОБЯЗАТЕЛЬНО!) Путь к файлу")
	outDir := flag.String("o", fmt.Sprintf("files_%s", time.Now().Format("2006-01-02_15-04-05")), "Out dir name")
	filterColumnNum := flag.Int("c", 5, "Номер колонки по которой будет фильтроваться таблица. Нумерация начинается с 0")
	flag.BoolFunc("v", "Версия", printVersion)
	flag.Parse()

	execPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("execPath:", execPath)
	if *fileName == "" {
		err := selectFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	data, err := ods.Read(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	var title []string
	var rows [][]string
	var filterColumn []string
	table0 := data.Content.Body.Spreadsheet.Table[0]
	for i, row := range table0.TableRow {
		if i == 0 {
			for _, cell := range row.TableCell {
				title = append(title, cell.P)
			}
			continue
		}
		if len(row.TableCell) < 2 {
			continue
		}
		var r []string
		for _, cell := range row.TableCell {
			r = append(r, cell.P)
		}
		rows = append(rows, r)
		if filterColumnValue := row.TableCell[*filterColumnNum].P; filterColumnValue != "" {
			filterColumn = append(filterColumn, filterColumnValue)
		} else {
			filterColumn = append(filterColumn, "nil")
		}

	}
	uniq := unique(filterColumn)

	err = os.MkdirAll(*outDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, filter := range uniq {
		var res [][]string
		for _, row := range rows {
			if row[*filterColumnNum] == filter || filter == "nil" && row[*filterColumnNum] == "" {
				res = append(res, row)
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			writeFile(res, title, *outDir, filter+".csv")

		}()
	}
	wg.Wait()

	log.Println(time.Since(t))
}

func unique(input []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func writeFile(rows [][]string, title []string, path, filename string) {
	var dataSplit []string
	dataSplit = append(dataSplit, strings.Join(title, ";"))

	for _, row := range rows {
		dataSplit = append(dataSplit, strings.Join(row, ";"))
	}

	fileNameReplace := strings.Replace(filename, "/", "-", 2)
	var filePath = path + "/" + fileNameReplace
	outFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			panic(err)
		}
	}(outFile)

	_, err = outFile.Write([]byte(strings.Join(dataSplit, "\n")))
	if err != nil {
		panic(err)
	}
}

func printVersion(s string) error {
	if s == "true" {
		info, _ := debug.ReadBuildInfo()
		fmt.Println(info.Main.Version)
		os.Exit(0)
	} else {
		log.Fatal("Undefined ", s)
	}
	return nil
}

func selectFile(file *string) error {
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Height(10).
				Picking(true).
				Title("Table").
				Description("Выберите таблицу в формате .ods").
				AllowedTypes([]string{".ods"}).
				Value(file),
		),
	).WithShowHelp(true).Run()
	if err != nil {
		return err
	}
	return nil
}
