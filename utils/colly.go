package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	baseUrl = "http://www.szlawyers.com"
	number  = 1
	c *colly.Collector
	writer *csv.Writer
)

func DoWork(query string) {

	fName := "lawFirmDetails.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer = csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"", "律所名称", "设立时间", "负责人", "律师人数", "联系电话", "详情链接","办公地址"})



	// Instantiate default collector
	c = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: 5,
		RandomDelay: 1 * time.Second,
	})

	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL.String())
	//})

	//c.OnResponse(func(r *colly.Response) {
	//	fmt.Println("response with code: ", r.StatusCode)
	//})

	// main page
	c.OnHTML(".tab_list tr:not(:first-child, :last-child)", func(e *colly.HTMLElement) {

		row := e.ChildTexts("td")
		link := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		row = append([]string{strconv.Itoa(number)}, row[0:]...)
		row = append(row, link)

		err := writer.Write(row)
		if err != nil{
			fmt.Println("信息写入错误！！！程序退出运行")
			os.Exit(1)
		}
		number += 1
	})

	// go next
	c.OnHTML(".page a:last-child", func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))
		text := e.DOM.Text()
		if text == "下一页"{
			c.Visit(link)
		}
	})

	// member detail page
	//c.OnHTML(".list table[style*=\"word-break\"] tbody tr:nth-child(4), .list table[style*=\"word-break\"] tbody tr:nth-child(5), .list table[style*=\"word-break\"] tbody tr:nth-child(11)", func(e *colly.HTMLElement) {
	//
	//	personalDetail = e.ChildTexts("td")
	//	personalDetail = personalDetail[1:]
	//	fmt.Println(personalDetail)
	//})

	url := baseUrl + query
	c.Visit(url)
	c.Wait()
	log.Printf("Scraping finished, check file %q for results\n", fName)
}

//func gotoCompanyDetails(row []string,link string) {
//
//	//link := strings.Join(row[6:],"")
//	//c.Visit(link)
//	b := colly.NewCollector()
//	b.OnHTML(".list tbody tr:nth-child(12)", func(e *colly.HTMLElement) {
//		location := e.ChildTexts("td")
//		location = location[1:]
//		row = append(row,location...)
//
//		err := writer.Write(row)
//		if err != nil{
//			fmt.Println("信息写入错误！！！程序退出运行")
//			os.Exit(1)
//		}
//		//link := e.Request.AbsoluteURL(e.ChildAttr(".list tbody tr:nth-child(11) span a","href"))
//	})
//
//	b.Visit(link)
//}
