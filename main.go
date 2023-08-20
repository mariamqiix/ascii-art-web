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
	/*
		http.HandleFunc("/submit-form", handler)
		if err != nil {
			log.Fatalln("There's an error with the server:", err)
		}
	*/
	http.Handle("/", http.FileServer(http.Dir("./template")))
	http.HandleFunc("/v", serveIndex)
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v" {
		fmt.Fprint(w, "page not found (cutom message)")
		return
	}

	var Words [][]string
	Text1 := strings.ReplaceAll(r.FormValue("thetext"), "\\t", "   ")

	for j := 0; j < len(Text1); j++ {
		Words = append(Words, ReadLetter(Text1[j], r.FormValue("chose")))
	}
	var b [8]string
	for x := 0; x < 8; x++ {
		for n := 0; n < len(Words); n++ {
			b[x] += Words[n][x]
		}
	}

	indexTemplate, _ := template.ParseFiles("template/index.html")

	if r.Method == http.MethodPost {
		err := indexTemplate.Execute(w, b)
		if err != nil {
			fmt.Print(err)
		}
		return
	}
	var c []string
	err := indexTemplate.Execute(w, c)
	if err != nil {
		fmt.Print(err)
	}
}

func serveForm(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		fmt.Fprint(w, "bad use")
		return
	}

	// var Words [][]string
	// Text1 := strings.ReplaceAll(r.FormValue("thetext"), "\\t", "   ")
	// for j := 0; j < len(Text1); j++ {
	// 	Words = append(Words, ReadLetter(Text1[j], r.FormValue("chose")))
	// }

	/*
		err = r.ParseForm()
		if err != nil {
			fmt.Print(w, "error 404")
			return
		}
	*/

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
