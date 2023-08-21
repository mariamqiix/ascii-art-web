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

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v" {
		fmt.Fprint(w, "page not found (cutom message)")
		return
	}
	text := r.FormValue("thetext")
	if len(text)>100 {
		fmt.Fprint(w, "لاتكثر كلام")
		return
	}
	_ , error := os.Stat(r.FormValue("chose")+".txt")
	// check if error is "file not exists"
	if os.IsNotExist(error) {
		fmt.Fprint(w,r.FormValue("chose")+"file does not exist\n")
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
	// for i := range b {
	// 	b[i] = html.UnescapeString(b[i])
	// }
	indexTemplate, _ := template.ParseFiles("template/index.html")
	if r.Method == http.MethodPost {
		err := indexTemplate.Execute(w, n)
		if err != nil {
			fmt.Print(err)
		}
		return
	}
	// var c []string
	// err := indexTemplate.Execute(w, c)
	// if err != nil {
	// 	fmt.Print(err)
	// }
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

// func serveForm(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		fmt.Fprint(w, "bad use")
// 		return
// 	}

// 	// var Words [][]string
// 	// Text1 := strings.ReplaceAll(r.FormValue("thetext"), "\\t", "   ")
// 	// for j := 0; j < len(Text1); j++ {
// 	// 	Words = append(Words, ReadLetter(Text1[j], r.FormValue("chose")))
// 	// }

// 	/*
// 		err = r.ParseForm()
// 		if err != nil {
// 			fmt.Print(w, "error 404")
// 			return
// 		}
// 	*/

// }
