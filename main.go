package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PageData struct {
	Color string
	Text  []string
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/style.css" {
		http.ServeFile(w, r, "./template/style.css")
		return
	} else if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "./template/404.html")
		return
	} else if r.Method == "GET" {
		indexTemplate, _ := template.ParseFiles("template/index.html")
		var v []string
		pageData := PageData{
			Text:  v,
			Color: "black",
		}
		indexTemplate.Execute(w, pageData)
		return
	} else if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "./template/400.html") // should be 400
		return
	} else if len(r.FormValue("thetext")) > 100 {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "template/400.html")
		return
	}
	_, error := os.Stat(r.FormValue("chose") + ".txt")
	if os.IsNotExist(error) {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "./template/500.html")
		return
	}
	if r.FormValue("thetext") == "" {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "template/400.html")
		return
	}
	if !CheckLetter(r.FormValue("thetext")) {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "template/400.html")
		return
	}
	if !CheckColor(r.FormValue("color")) {
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "./template/400.html")
		return
	}

	TextInASCII := serveIndex(r.FormValue("thetext"), r.FormValue("chose"))
	indexTemplate, _ := template.ParseFiles("template/index.html")
	if r.Method == "POST" {
		color := "black"
		if r.FormValue("color") != "#f1f0e8" {
			color = r.FormValue("color")
		}
		pageData := PageData{
			Text:  TextInASCII,
			Color: color,
		}
		err := indexTemplate.Execute(w, pageData)
		if err != nil {
			fmt.Print(err)
		}
		return
	}
}

func serveIndex(text, filename string) []string {
	WordsInArr := strings.Split(text, "\r\n")
	var Text []string
	for l := 0; l < len(WordsInArr); l++ {
		var Words [][]string
		Text1 := strings.ReplaceAll(WordsInArr[l], "\\t", "   ")
		if Text1 != "" {
			for j := 0; j < len(Text1); j++ {
				Words = append(Words, ReadLetter(Text1[j], filename))
			}
			for x := 0; x < 8; x++ {
				Lines := ""
				for n := 0; n < len(Words); n++ {
					Lines += Words[n][x]
				}
				Text = append(Text, Lines)
			}
		} else {
			Text = append(Text, "\n")
		}
	}
	Line := strings.Join(Text, "\n")
	var Words []string
	Words = append(Words, Line)
	return Words
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

func CheckLetter(s string) bool {
	WordsInArr := strings.Split(s, "\r\n")
	for l := 0; l < len(WordsInArr); l++ {
		for g := 0; g < len(WordsInArr[l]); g++ {
			if WordsInArr[l][g] > 126 || WordsInArr[l][g] < 32 {
				return false
			}
		}
	}
	return true
}

func CheckColor(userValue string) bool {
	Colors := []string{"red", "green", "yellow", "blue", "purple", "cyan", "white", "black", "orange"}
	for _, color := range Colors {
		if color == userValue {
			return true
		} else if strings.Index(userValue, "rgb") == 0 && userValue[len(userValue)-1] == ')' {
			userValue := strings.ReplaceAll(userValue, " ", "")
			c := (strings.Split(userValue[4:len(userValue)-1], ","))
			for i := 0; i < len(c); i++ {
				check, err := strconv.Atoi(c[i])
				if err != nil || len(c) != 3 || check < 0 || check > 255 {
					return false
				}
			}
			return false
		} else if strings.Index(userValue, "#") == 0 && len(userValue) == 7 {
			for i := 1; i <= 6; i++ {
				if (userValue[i] >= '0' && userValue[i] <= '9') || (userValue[i] >= 'a' && userValue[i] <= 'f') {
				} else {
					return false
				}
			}
			return true
		}
	}
	return false
}
