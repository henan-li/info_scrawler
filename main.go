package main

import (
	"fmt"
	"html/template"
	"info_scrawler/utils"
	"info_scrawler/view"
	"log"
	"net/http"
	_ "strings"
)

func main() {

	dirs := []string{"view"} // 设置需要释放的目录
	for _, dir := range dirs {
		// 解压dir目录到当前目录
		if err := view.RestoreAssets("./", dir); err != nil {
			break
		}
	}

	http.HandleFunc("/", indexPage)    // display form
	http.HandleFunc("/info", infoPage) // handle get request

	fmt.Println("server will start and listen at localhost:8080\nuse url: localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 展示表单
func indexPage(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseFiles("./view/index.html"))
	t.Execute(w, "")
}

// 处理表单参数
func infoPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var str = "/searchLawFirm?"
	for k, v := range query {
		str += k + "=" + v[0] + "&"
	}

	var firmType string
	if query["lawFirmType"][0] == "235bbe7b44ea4eb381e816e7436f8afa"{
		firmType = "personal"
	}else if query["lawFirmType"][0] == "592f86a85a9b4b3d98e24db19f3ae93b"{
		firmType = "group"
	}else{
		firmType = ""
	}
	// write content into .csv file
	utils.DoWork(str,firmType)
}
