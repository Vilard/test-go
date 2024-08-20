package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func createReport() {
	date := time.Now()
	strDate := date.Format("02-01-2006 15:04:05")
	text := "тут отчет:" + strDate
	filename := strDate + ".log"
	file, err := os.Create(filename)

	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer file.Close()
	file.WriteString(text)

	fmt.Println("File created........" + filename)

	var buff bytes.Buffer

	zipW := zip.NewWriter(&buff)
	f, err := zipW.Create(filename)

	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(text))
	if err != nil {
		panic(err)
	}

	err = zipW.Close()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./ui/static/report.zip", buff.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func crHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodPost {
		createReport()
		fmt.Fprintf(w, "Report created successfully!")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Control-Allow-Credentials", "true")

	http.ServeFile(w, r, "static/report.zip")
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/create-report", crHandler)
	mux.HandleFunc("/download", downloadHandler)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
