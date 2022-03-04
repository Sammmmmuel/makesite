package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"time"

)

type Article struct {
	Title       string
	Paragraphs  []string
	ReadingTime string
}

func parseFile(filePath string) Article {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	article := Article{}

	for index, line := range strings.Split(string(fileContents), "\n") {
		if index == 0 {
			article.Title = line
		} else if line != "" {
			article.Paragraphs = append(article.Paragraphs, line)
		}
	}

	return article
}

func generateHtml(article Article, fileName string) int64 {
	t := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))

	output, err := os.Create("dist/" + fileName + ".html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(output, article)
	if err != nil {
		panic(err)
	}

	info, err := os.Stat("dist/" + fileName + ".html")
	if err != nil {
		panic(err)
	}

	return info.Size()
}

func main() {
	start := time.Now()
	var generatedSize int64 = 0

	filePath := flag.String("file", "", "Parse file from the given path.")
	dirPath := flag.String("dir", "", "Parse all files from the given directory.")
	flag.Parse()
	filesGenerated := 0

	if *filePath != "" {
		splitPath := strings.Split(*filePath, "/")
		fileName := splitPath[len(splitPath)-1]
		fileName = strings.Split(fileName, ".")[0]

		article := parseFile(*filePath)
		generatedSize += generateHtml(article, fileName)

		filesGenerated += 1
	}

	if *dirPath != "" {
		files, err := ioutil.ReadDir(*dirPath)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".txt") {
				fileName := strings.Split((file.Name()), ".")[0]

				article := parseFile(*dirPath + "/" + file.Name())
				generatedSize += generateHtml(article, fileName)

				filesGenerated += 1
			}
		}
	}

	elapsed := time.Since(start)
	kilobytes := float64(generatedSize) / 1024

	fmt.Print("Success!")
	fmt.Println(" Generated " + fmt.Sprintf("%d", filesGenerated) + " pages " + "(" + fmt.Sprintf("%.2f", kilobytes) + "kB total)" + " in " + fmt.Sprintf("%.3f", elapsed.Seconds()) + " seconds.")
}

