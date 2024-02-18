package main

import (
	"md2htm/utils"
	"os"
	"strings"
)

func main() {
	fileName := os.Args[1]
	fileNameWithoutExtension := strings.Split(fileName, ".")[0]
	extension := strings.Split(fileName, ".")[1]
	if extension != "md" && extension != "markdown" {
		panic("File type not supported")
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(file))
	fileContents := string(file)
	fileLines := strings.Split(fileContents, "\n")
	if fileLines[len(fileLines)-1] == "" {
		fileLines = fileLines[:len(fileLines)-1]
	}
	for i := 0; i < len(fileLines); i++ {
		// check if the prefix is in the map
		for k := range utils.Tags {
			if strings.HasPrefix(fileLines[i], k) {
				if k == "- " || k == "* " {
					if i == 0 || (!strings.HasPrefix(fileLines[i-1], "<ul>")) {
						fileLines[i] = "<ul><li>" + fileLines[i][2:] + "</li>"
					} else {
						fileLines[i] = "<li>" + fileLines[i][2:] + "</li>"
					}
					if i == len(fileLines)-1 || (!strings.HasPrefix(fileLines[i+1], "- ") && !strings.HasPrefix(fileLines[i+1], "* ")) {
						fileLines[i] = fileLines[i] + "</ul>"
					}
				} else {
					if i > 0 && (strings.HasPrefix(fileLines[i-1], "- ") || strings.HasPrefix(fileLines[i-1], "* ")) {
						fileLines[i-1] = fileLines[i-1] + "</ul>"
					}
					fileLines[i] = utils.ConvertToHTMLTags(k, fileLines[i])
				}
			}
		}
	}
	htmlFormatOfFile := strings.Join(fileLines, "\n")
	templateHtml, err := os.ReadFile("template.html")
	if err != nil {
		panic(err)
	}
	// convert fileNameWithoutExtension to title case
	templateHtml = []byte(strings.Replace(string(templateHtml), "$data", htmlFormatOfFile, 1))
	os.WriteFile(fileNameWithoutExtension+".html", templateHtml, 0644)
}
