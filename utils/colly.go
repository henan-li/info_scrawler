package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

var(
	baseUrl = "http://www.szlawyers.com"
)

func DoWork(query string) {
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
	writer.Write([]string{"律所名称", "设立时间", "负责人", "律师人数", "联系电话", "详情链接"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("tbody tr:not(:first-child, :last-child)", func(e *colly.HTMLElement) {

		// 流程: (注意:他是一次一个的去找,因此每一个row就对应页面上的一个row) 从tbody tr下面找td,然后将他的文本内容都取出来
		// 然后在往当前row中动态追加一个url
		// 最后写入csv,然后开始下一轮查找
		row := e.ChildTexts("td")
		link := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		row = append(row,link)
		writer.Write(row)
	})

	url := baseUrl+query
	fmt.Println(url)
	c.Visit(url)

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
