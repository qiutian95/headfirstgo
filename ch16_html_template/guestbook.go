package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type GuestBook struct {
	SignatureCount int
	Signatures     []string
}

// 查看签名主页
func viewHandler(w http.ResponseWriter, r *http.Request) {
	signatures := getStrings("signatures.txt") // 这里一开始老是找不到，需要调整build config中的working directory,调整为当前项目目录
	html, err := template.ParseFiles("view.html")
	check(err)
	guestBook := GuestBook{
		SignatureCount: len(signatures),
		Signatures:     signatures,
	}
	err = html.Execute(w, guestBook)
	check(err)
}

// 新增签名页面
func newHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("new.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)
}

// 提交签名，将签名写入文件
func createHandler(w http.ResponseWriter, r *http.Request) {
	signature := r.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("signatures.txt", options, os.FileMode(0600))
	check(err)
	_, err = fmt.Fprintln(file, signature)
	check(err)
	err = file.Close() // 这里是写文件，需要实时处理错误
	check(err)
	//viewHandler(w, r) // 这种偷懒的方式是不是也可以（这种的话url显示的是当前的url）

	//重定向更好
	http.Redirect(w, r, "/guestbook", http.StatusFound)

}

func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) { // 文件不存在返回nil
		return nil
	}
	check(err)
	defer func(file *os.File) {
		err := file.Close()
		check(err)
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/new", newHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
