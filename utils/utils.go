package utils

import "strings"

var Tags = map[string]string{
	"# ":   "h1",
	"## ":  "h2",
	"### ": "h3",
	"- ":   "li",
	"* ":   "li",
	"** ":  "strong",
	"__ ":  "strong",
	"~~ ":  "del",
	"` ":   "code",
}

var HtmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="author" content="$authorName">
  <meta name="description" content="$description">
  <title>$pageTitle</title>
</head>
<body>
  $data
</body>
</html>
`

var Metadata = []string{"pageTitle", "authorName", "description"}

func ConvertToHTMLTags(mdPrefix string, lineContent string) string {
	tag, ok := Tags[mdPrefix]
	if lineContent == "---" {
		return ""
	}
	lineContent = lineContent[len(mdPrefix):]
	if !ok {
		return "<p>" + lineContent + "</p>"
	}
	return "<" + tag + ">" + lineContent + "</" + tag + ">"
}

func HandleMetadata(fileLines *[]string, metadataValues *map[string]string, i int) {
	if i == 0 && (*fileLines)[i] == "---" {
		for j := i + 1; j < len(*fileLines); j++ {
			if (*fileLines)[j] == "---" {
				*fileLines = append((*fileLines)[:i], (*fileLines)[j+1:]...)
				break
			}
			for _, v := range Metadata {
				if strings.HasPrefix((*fileLines)[j], v) {
					(*metadataValues)[v] = (*fileLines)[j][len(v)+2:]
				}
			}
		}
	}
}

func HandleLists(fileLines *[]string, i int, k string) {
	lines := *fileLines
	if i == 0 || (!strings.HasPrefix(lines[i-1], "<ul>")) {
		lines[i] = "<ul><li>" + lines[i][2:] + "</li>"
	} else {
		lines[i] = "<li>" + lines[i][2:] + "</li>"
	}
	if i == len(lines)-1 || (!strings.HasPrefix(lines[i+1], "- ") && !strings.HasPrefix(lines[i+1], "* ")) {
		lines[i] = lines[i] + "</ul>"
	}
}
