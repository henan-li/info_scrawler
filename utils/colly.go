package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func main() {
	fName := "lawFirmDetails.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"律所名称", "设立时间", "负责人", "律师人数", "联系电话"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".tab_list tbody tr", func(e *colly.HTMLElement) {

		res := e.ChildTexts("tr[bgcolor=\"#f0f0f0\"] td")
		fmt.Println(len(res))
		writer.Write(res)
	})

	c.Visit("http://www.szlawyers.com/searchLawFirm?name=&creditCode=&justiceBureauId=&officeZone=&beginPracticeLicenseDate=&endPracticeLicenseDate=&lawFirmType=&scale=100&x=27&y=13")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}