package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	Output [][]string
}

var tmplt *template.Template

func main() {
	serverfile := http.FileServer(http.Dir("./template"))
	http.Handle("/", serverfile)
	fmt.Println("Starting server....")
	http.HandleFunc("/submit-form", handler)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatalln("There's an error with the server:", err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	var Words [][]string
	Text1 := strings.ReplaceAll( r.FormValue("thetext"), "\\t", "   ")
	for j := 0; j < len(Text1); j++ {
		Words = append(Words, ReadLetter(Text1[j], r.FormValue("chose")))
	}
	_, err := template.New("foo").Parse("Hello")
	err = r.ParseForm()
	if err != nil {
		fmt.Print(w, "error 404")
		return
	}

	for x := 0; x < 8; x++ {
		for n := 0; n < len(Words); n++ {
			fmt.Fprintf(w, Words[n][x])
		}
		if x+1 != 8 {
			fmt.Fprintf(w, "\n")
		}
	}
	fmt.Fprintf(w, "\n")
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
