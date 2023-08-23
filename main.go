package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./template")))
	http.HandleFunc("/v", serveIndex)
	http.ListenAndServe(":8080", nil)
}

// func Handler(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		http.ServeFile(w, r, "./template/404.html")
// 	} else {
// 		http.ServeFile(w, r, "./template/index.html")
// 	}
// }

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v" && r.URL.Path != "/" {
		http.ServeFile(w, r, "./template/404.html")
		return
	}

	// if r.URL.Path == "/" {
	// 	http.ServeFile(w, r, "./template/index.html")
	// } else {
	// 	http.ServeFile(w, r, "./template/404.html")
	// }

	text := r.FormValue("thetext")
	if len(text) > 100 {
		http.ServeFile(w, r, "template/500.html")
		return
	}
	_, error := os.Stat(r.FormValue("chose") + ".txt")
	// check if error is "file not exists"
	if os.IsNotExist(error) {
		http.ServeFile(w, r, "template/404.html")
		return
	}
	WordsInArr := strings.Split(text, "\r\n")

	var b []string
	for l := 0; l < len(WordsInArr); l++ {
		var Words [][]string
		Text1 := strings.ReplaceAll(WordsInArr[l], "\\t", "   ")
		for j := 0; j < len(Text1); j++ {
			Words = append(Words, ReadLetter(Text1[j], r.FormValue("chose")))
		}
		for x := 0; x < 8; x++ {
			Lines := ""
			for n := 0; n < len(Words); n++ {
				Lines += Words[n][x]
			}
			b = append(b, Lines)
		}
	}
	newB := strings.Join(b, "\n")
	var n []string
	n = append(n, newB)
	indexTemplate, _ := template.ParseFiles("template/index.html")
	if r.Method == "POST"{
		err := indexTemplate.Execute(w, n)
		if err != nil {
			fmt.Print(err)
		}
		return
	} else {
		http.ServeFile(w, r, "template/400.html")
		return
	}
}

func ReadLetter(Text1 byte, fileName string) []string {
	//buffio object, to open and read the file
	ReadFile, _ := os.Open(fileName + ".txt")
	FileScanner := bufio.NewScanner(ReadFile)
	var Letter []string
	stop := 1
	i := 0
	a := (int(Text1)-32)*9 + 2
	for FileScanner.Scan() {
		i++
		if i >= a {
			stop++
			Letter = append(Letter, FileScanner.Text())
			if stop > 8 {
				break
			}
		}
	}
	ReadFile.Close()
	return Letter
}
