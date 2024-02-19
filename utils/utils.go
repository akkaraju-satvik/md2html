package utils

import (
	"strings"
)

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
	"```":  "pre",
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
	<style>
		* {
			font-family: sans-serif;
		}
		pre {
			width: max-content;
			padding: 1em;
			background-color: #e0e0e0;
			border-radius: 5px;
		}
		code {
			font-family: monospace;
		}
	</style>
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

func HandleMetadata(fileLines []string, metadataValues *map[string]string) int {
	var j int
	if fileLines[0] == "---" {
	x:
		for j = 1; j < len(fileLines); j++ {
			if (fileLines)[j] == "---" {
				break
			}
			for _, v := range Metadata {
				if strings.HasPrefix((fileLines)[j], v) {
					(*metadataValues)[v] = (fileLines)[j][len(v)+2:]
					continue x
				}
			}
		}
	}
	return j
}

func HandleLists(fileLines *[]string, i int, k string) int {
	lines := *fileLines
	if i == 0 || (!strings.HasPrefix(lines[i-1], "<ul>")) {
		lines[i] = "<ul><li>" + lines[i][2:] + "</li>"
	} else {
		lines[i] = "<li>" + lines[i][2:] + "</li>"
	}
	if i == len(lines)-1 || (!strings.HasPrefix(lines[i+1], "- ") && !strings.HasPrefix(lines[i+1], "* ")) {
		lines[i] = lines[i] + "</ul>"
	}
	return i
}

func HandleCodeBlocks(fileLines *[]string, i int) {
	var codeBlock, codeTag string
	language := strings.Split((*fileLines)[i], "```")[1]
	if language != "" {
		codeBlock = language + "\n"
		codeTag = "<code class=\"language-" + language + "\">"
	} else {
		codeTag = "<code>"
	}
	var j int
	for j = i; j < len((*fileLines)); j++ {
		if (*fileLines)[j] == "```" {
			codeBlock = strings.Join((*fileLines)[i+1:j], "\n")
			break
		}
	}
	(*fileLines)[i] = "<pre>" + codeTag + codeBlock + "</code></pre>"
	(*fileLines)[i] = strings.TrimSpace((*fileLines)[i])
	if j != len((*fileLines))-1 {
		(*fileLines) = append((*fileLines)[:i+1], (*fileLines)[j+1:]...)
	} else {
		(*fileLines) = (*fileLines)[:i+1]
	}
}
