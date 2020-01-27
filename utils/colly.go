package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

var (
	baseUrl     = "http://www.szlawyers.com"
	number      = 1
	//number2     = 1
	c           *colly.Collector
	writer      *csv.Writer
	writer2     *csv.Writer
	companyPage map[int][]string
)

func DoWork(query string,firmType string) {

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
	writer2.Write([]string{"", "律师姓名", "律师性别", "所属律所", "取得律师资格证时间", "在深圳开始执业时间", "在深圳执业时长(天)", "证件照"})

	// Instantiate default collector
	c = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		//colly.Async(true),
	)
	//c.Limit(&colly.LimitRule{
	//	Parallelism: 2,
	//	RandomDelay: 2 * time.Second,
	//})

	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL.String())
	//})

	//c.OnResponse(func(r *colly.Response) {
	//	fmt.Println("response with code: ", r.StatusCode)
	//})


	var addressSelector string
	var personalSelector string
	if firmType == "personal"{
		addressSelector = ".lawyer_info tbody tr:nth-child(11)"
		personalSelector = ".lawyer_info tbody tr:nth-child(10)"
	}else{
		addressSelector = ".lawyer_info tbody tr:nth-child(12)"
		personalSelector = ".lawyer_info tbody tr:nth-child(11)"
	}


	c.OnHTML(".tab_list tr:not(:first-child, :last-child)", func(e *colly.HTMLElement) {

		// main page
		row := e.ChildTexts("td")
		link := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		row = append([]string{strconv.Itoa(number)}, row[0:]...)
		row = append(row, link)

		// visit company address and get workplace info
		e.Request.Visit(link)
		row = append(row, companyPage[1]...)

		fmt.Println("working on company information: ", row)
		err := writer.Write(row)
		if err != nil {
			fmt.Println("公司信息写入错误！！！程序退出运行")
			os.Exit(1)
		}

		number += 1
	})

	// go next
	c.OnHTML(".page a:last-child", func(e *colly.HTMLElement) {

		link := e.Request.AbsoluteURL(e.Attr("href"))
		text := e.Text
		if text == "下一页" {
			e.Request.Visit(link)
		}
	})

	// visit personal page
	c.OnHTML(personalSelector, func(e *colly.HTMLElement) {

		urls := e.ChildAttrs("td:nth-child(2) span a", "href")

		//fmt.Println(len(test),test)
		for _, v := range urls {

			url := baseUrl + v
			flag, _ := e.Request.HasVisited(url)

			if !flag {
				e.Request.Visit(url)
			}
		}
		//urlLists := make(map[int]string)
		//
		//e.ForEach("td:nth-child(2) span", func(i int, element *colly.HTMLElement) {
		//	//urlLists[i] = element.Request.AbsoluteURL(element.ChildAttr("a", "href"))
		//	url := element.Request.AbsoluteURL(element.ChildAttr("a", "href"))
		//
		//})
		//
		//fmt.Println(urlLists)
		//for _, v := range urlLists {
		//	flag,_ :=e.Request.HasVisited(v)
		//	if !flag {
		//		e.Request.Visit(v)
		//	}
		//}

	})

	// personal page details tr:nth-child(2) td:nth-child(2) span 2 4 5 8 11
	c.OnHTML(".list table[style*=\"word-break\"] tbody", func(e *colly.HTMLElement) {

		personalRowRes := []string{""}
		personalRowRes = append(personalRowRes, e.ChildText("tr:nth-child(2) td:nth-child(2) span span"))
		personalRowRes = append(personalRowRes, e.ChildText("tr:nth-child(4) td:nth-child(2) span span"))
		personalRowRes = append(personalRowRes, e.ChildText("tr:nth-child(5) td:nth-child(2) span a"))
		personalRowRes = append(personalRowRes, e.ChildText("tr:nth-child(8) td:nth-child(2) span span"))

		startTime := e.ChildText("tr:nth-child(11) td:nth-child(2) span span")
		personalRowRes = append(personalRowRes, startTime)

		if startTime != ""{
			dayDiffRes := getDayDiff(startTime)
			personalRowRes = append(personalRowRes, dayDiffRes)
			personalRowRes = append(personalRowRes, e.Request.AbsoluteURL(e.ChildAttr("tr:nth-child(2) td:nth-child(3) img", "src")))

		}

		fmt.Println("working on personal details: ", personalRowRes)
		err := writer2.Write(personalRowRes)
		if err != nil {
			fmt.Println("个人信息写入错误！！！程序退出运行")
			os.Exit(1)
		}
	})

	// company page details

	c.OnHTML(addressSelector, func(e *colly.HTMLElement) {

		companyPage = make(map[int][]string)
		num := 1
		location := e.ChildTexts("td")
		location = location[1:]
		companyPage[num] = location
	})

	url := baseUrl + query
	c.Visit(url)
	//c.Wait()
	log.Printf("Scraping finished, check file %q for results\n", fName)
	log.Printf("Scraping finished, check file %q for results\n", fName2)
}

func getDayDiff(startTime string) string {

	// 移除中文字符：年月日， 重新组合成 yyyy-mm-dd
	sT := []byte(startTime)
	year := string(sT[0:4])
	month := string(sT[7:9])
	day := string(sT[12:14])
	startTimeFormat := year + "-" + month + "-" + day

	a, _ := time.Parse("2006-01-02", startTimeFormat)

	currentTime := time.Now().Format("2006-01-02")
	b, _ := time.Parse("2006-01-02", currentTime)

	dayDiff := b.Sub(a).Hours() / 24
	dayDiff = math.Ceil(dayDiff)

	dayDiffRes := strconv.FormatFloat(dayDiff, 'f', 0, 64)
	return dayDiffRes
}
