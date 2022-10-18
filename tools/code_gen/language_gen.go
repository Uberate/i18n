//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type LanguageMetadata struct {
	Custom   string
	ISO6391  string
	ISO6392T string
	ISO6392B string
}

var languages = map[string]LanguageMetadata{
	"ChineseLn":  {Custom: "chinese", ISO6391: "zh", ISO6392T: "zho", ISO6392B: "chi"},
	"EnglishLn":  {Custom: "english", ISO6391: "en", ISO6392B: "en", ISO6392T: "en"},
	"JapaneseLn": {Custom: "japanese", ISO6391: "ja", ISO6392B: "jap", ISO6392T: "jap"},
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error params")
		os.Exit(-1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	tmp, err := template.ParseFiles(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	bbf := &bytes.Buffer{}
	err = tmp.Execute(bbf, languages)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err = os.WriteFile(outputFile, bbf.Bytes(), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
