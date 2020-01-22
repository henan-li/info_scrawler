package main

import (
	"fmt"
	"html/template"
	"info_scrawler/utils"
	"log"
	"net/http"
	_ "strings"
)

func main() {
	// 请求文件路径转成服务器路径
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("view/index.html"))))
	// 路由处理
	http.HandleFunc("/", indexPage)    // display form
	http.HandleFunc("/info", infoPage) // handle get request

	// 启动服务
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

	// write content into .csv file
	utils.DoWork(str)
}

