package main

import (
	"fmt"
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
	x := string(file)

	xa := strings.Split(x, "\n")
	if xa[len(xa)-1] == "" {
		xa = xa[:len(xa)-1]
	}
	for i := 0; i < len(xa); i++ {
		// check if the prefix is in the map
		for k := range utils.Tags {
			if strings.HasPrefix(xa[i], k) {
				if k == "- " || k == "* " {
					if i == 0 || (!strings.HasPrefix(xa[i-1], "<ul>")) {
						xa[i] = "<ul><li>" + xa[i][2:] + "</li>"
					} else {
						xa[i] = "<li>" + xa[i][2:] + "</li>"
					}
					if i == len(xa)-1 || (!strings.HasPrefix(xa[i+1], "- ") && !strings.HasPrefix(xa[i+1], "* ")) {
						xa[i] = xa[i] + "</ul>"
					}
				} else {
					if i > 0 && (strings.HasPrefix(xa[i-1], "- ") || strings.HasPrefix(xa[i-1], "* ")) {
						xa[i-1] = xa[i-1] + "</ul>"
					}
					xa[i] = utils.ConvertToHTMLTags(k, xa[i])
				}
			}
		}
	}
	html := strings.Join(xa, "\n")
	fmt.Println(html)
	templateHtml, err := os.ReadFile("template/template.html")
	if err != nil {
		panic(err)
	}
	os.WriteFile(fileNameWithoutExtension+".html", []byte(strings.Replace(string(templateHtml), "$data", html, 1)), 0644)
}
