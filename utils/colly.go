package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
)

var (
	baseUrl      = "http://www.szlawyers.com"
	number       = 1
	number2      = 1
	c            *colly.Collector
	writer       *csv.Writer
	writer2      *csv.Writer
	companyPage  map[int][]string
	personalPage map[int][]string
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
	writer.Write([]string{"", "律所名称", "设立时间", "负责人", "律师人数", "联系电话", "详情链接", "办公地址"})

	fName2 := "lawFirmPersonDetails.csv"
	file2, err := os.Create(fName2)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName2, err)
		return
	}
	defer file2.Close()
	writer2 = csv.NewWriter(file2)
	defer writer2.Flush()
	writer2.Write([]string{"", "律师姓名", "律师性别", "所属律所", "取得律师资格证时间", "在深圳开始执业时间"})

	// Instantiate default collector
	c = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		//colly.Async(true),
	)
	//c.Limit(&colly.LimitRule{
	//	Parallelism: 5,
	//	RandomDelay: 1 * time.Second,
	//})

	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL.String())
	//})

	//c.OnResponse(func(r *colly.Response) {
	//	fmt.Println("response with code: ", r.StatusCode)
	//})

	c.OnHTML(".tab_list tr:not(:first-child, :last-child)", func(e *colly.HTMLElement) {

		// main page
		row := e.ChildTexts("td")
		link := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		row = append([]string{strconv.Itoa(number)}, row[0:]...)
		row = append(row, link)

		// visit company address and get workplace info
		visitSubPage(link)
		row = append(row, companyPage[1]...)

		err := writer.Write(row)
		if err != nil {
			fmt.Println("信息写入错误！！！程序退出运行")
			os.Exit(1)
		}

		number += 1
	})

	// go next
	c.OnHTML(".page a:last-child", func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))
		text := e.DOM.Text()
		if text == "下一页" {
			c.Visit(link)
		}
	})

	// company page details
	c.OnHTML(".lawyer_info tbody tr:nth-child(12)", func(e *colly.HTMLElement) {

		companyPage = make(map[int][]string)
		num := 1
		location := e.ChildTexts("td")
		location = location[1:]
		companyPage[num] = location
	})

	// visit personal page
	c.OnHTML(".lawyer_info tbody tr:nth-child(11)", func(e *colly.HTMLElement) {

		nameList := make(map[int]string)
		e.ForEach("td span", func(i int, element *colly.HTMLElement) {
			nameList[i] = element.Request.AbsoluteURL(element.ChildAttr("a", "href"))
		})
		//nameList := e.Request.AbsoluteURL(e.ChildAttr("td:nth-child(2) span a","href"))

		for _, v := range nameList {
			visitSubPage(v)
		}
	})

	// personal page details
	c.OnHTML(
		".list table[style*=\"word-break\"] tbody", func(e *colly.HTMLElement) {

			//err := writer2.Write(personalPage[1])
			//if err != nil {
			//	fmt.Println("信息写入错误！！！程序退出运行")
			//	os.Exit(1)
			//}
			//
			//number2 += 1
			fmt.Println(personalPage)
		})

	url := baseUrl + query
	c.Visit(url)
	c.Wait()
	log.Printf("Scraping finished, check file %q for results\n", fName)
}

func visitSubPage(link string) {

	//link := strings.Join(row[6:],"")
	//c.Visit(link)
	flag, _ := c.HasVisited(link)
	if !flag {
		c.Visit(link)
	}
}
