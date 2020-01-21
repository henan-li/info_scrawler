package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	_ "strings"
)

func main() {

	// 路由处理
	http.HandleFunc("/", indexPage)    // display form
	http.HandleFunc("/info", infoPage) // handle get request

	// 启动服务
	fmt.Println("server will start and listen at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 展示表单
func indexPage(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseFiles("view/index.html"))
	t.Execute(w, "")
}

// 处理表单参数
func infoPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var str = ""
	for k, v := range query {
		str += k + "=" + v[0] + "&"
	}

	// remove last &
	//str = strings.TrimRight(str, "&")

	// write content into .csv file
	pageContent := getContent(str)
	fmt.Println(pageContent)

}

// 根据参数获取所有数据
func getContent(str string) []string {

	//var page string

	str = "http://www.szlawyers.com/searchLawFirm?" + str
	response, err := http.Get(str)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	// all page info
	content := string(contents)

	// get title
	var allTitle []string
	reg := regexp.MustCompile(`(律所名称)|(设立时间)|(负责人)|(律师人数)|(联系电话)`)
	titles := reg.FindAllStringSubmatch(content, -1)
	fmt.Println(titles)
	for _, v := range titles {
		allTitle = append(allTitle, v[2])
	}
	fmt.Println(allTitle)
	//allTitle = allTitle[2:]

	// get content
	//var allContent []UserInfo
	//reg = regexp.MustCompile(``)
	//lists := reg.FindAllStringSubmatch(content, -1)
	//fmt.Println(lists)

	//UserInfo := &UserInfo{}
	//for _, v := range lists {
	//	allContent = append(allContent, v[0])
	//}

	return allTitle
}

type UserInfo struct {
	FirmName      string
	Time          string
	Principal     string
	Lawyers       string
	ContactNumber string
}
