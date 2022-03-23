package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/mrizkyy46/XenditTest/model"
)

func main() {
	headers := []string{
		"Amt",
		"Descr",
		"Date",
		"ID",
		"Remarks",
	}

	//assign type of struct to variable
	var proxyList []model.Proxy
	var sourceList []model.Source
	var exportList []model.Export

	//read all CSV lines and append to struct
	sourceLines := readAllCSV("export/source.csv")

	for _, i := range sourceLines {
		sourceList = append(sourceList, model.Source{
			Date: i[0], ID: i[1], Amount: i[2], Description: i[3],
		})
	}
	proxyLines := readAllCSV("export/proxy.csv")

	for _, i := range proxyLines {
		proxyList = append(proxyList, model.Proxy{
			Amt: i[0], Descr: i[1], Date: i[2], ID: i[3],
		})
	}

	sort.Slice(proxyList, func(i, j int) bool {
		return proxyList[i].Date < proxyList[j].Date
	})

	//get proxy data only for july 2021
	for key := range proxyList {
		proxyDate := substr(proxyList[key].Date, 5, 7)
		remarks := ""

		if proxyDate == "07" {
			resSourceID := findID(proxyList[key].ID, sourceList)

			if !resSourceID {
				remarks = "DISCREPANCY DATA - This is non-reconciled entries"
			}

			exportList = append(exportList, model.Export{
				Amt:     proxyList[key].Amt,
				Descr:   proxyList[key].Descr,
				Date:    proxyList[key].Date,
				ID:      proxyList[key].ID,
				Remarks: remarks,
			})

		}
	}

	//---------------------------------------------------------------------------------------------

	//create new csv file
	file, err := os.Create("export/export.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	writer := csv.NewWriter(file)

	defer writer.Flush()
	writer.Write(headers)

	for key := range exportList {
		r := make([]string, 0, 1+len(headers))

		r = append(
			r,
			exportList[key].Amt,
			exportList[key].Descr,
			exportList[key].Date,
			exportList[key].ID,
			exportList[key].Remarks,
		)

		writer.Write(r)
	}
	writer.Flush()

	//---------------------------------------------------------------------------------------------

	textFile, err := os.Create("export/summary.txt")
	if err != nil {
		fmt.Println(err)
		textFile.Close()
		return
	}

	countProcess := len(exportList)
	var countDiscrepancies int
	var typeDiscrepancies string

	for key := range exportList {
		if exportList[key].Remarks != "" {
			countDiscrepancies += 1
			typeDiscrepancies = exportList[key].Remarks
		}
	}

	d := []string{
		"Summary Report",
		"",
		"Date range			: " + exportList[0].Date + " - " + exportList[countProcess-1].Date,
		"Source records Processed	: " + strconv.Itoa(countProcess),
		"Number of discrepancies		: " + strconv.Itoa(countDiscrepancies),
		"Type of discrepancies		: " + typeDiscrepancies,
	}

	for _, v := range d {
		fmt.Fprintln(textFile, v)
	}

	fmt.Println("file written successfully")
}

func readAllCSV(filepath string) [][]string {

	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	datas := [][]string{}
	readLines := csv.NewReader(file)

	for {
		line, err := readLines.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		datas = append(datas, line)
	}

	return datas[1:]
}

func substr(s string, start, end int) string {
	counter, startIdx := 0, 0
	for i := range s {
		if counter == start {
			startIdx = i
		}
		if counter == end {
			return s[startIdx:i]
		}
		counter++
	}
	return s[startIdx:]
}

func findID(id string, sources []model.Source) (result bool) {
	result = false
	for _, proxy := range sources {
		if proxy.ID == id {
			result = true
			break
		}

	}

	return result
}
