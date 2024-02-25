package utils

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var Tags = map[string]string{
	"#":   "h1",
	"##":  "h2",
	"###": "h3",
	"-":   "li",
	"*":   "li",
	"**":  "strong",
	"__":  "strong",
	"~~":  "del",
	"`":   "pre",
	"```": "pre",
}

var HtmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="author" content="$authorName &lt;$authorEmail&gt;">
  <meta name="description" content="$description">
	<link rel="icon" href="$favicon" type="image/x-icon">
  <title>$pageTitle | $projectName</title>
	<style>
		* {
			font-family: sans-serif;
		}
		pre {
			width: max-content;
			padding: 1em;
			display: inline-block;
			background-color: #e0e0e0;
			border-radius: 5px;
		}
		code {
			font-family: monospace;
			background-color: #e0e0e0;
			padding: 0.2em;
		}
		.container {
			width: 80%;
			margin: 0 auto;
		}
		h1 {
			text-align: center;
		}
		a {
			color: #000000;
		}
	</style>
</head>
<body>
	<div class="container">
  	$data
	</div>
</body>
</html>
`

type Config struct {
	ProjectName string `yaml:"projectName"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Favicon     string `yaml:"favicon"`
	Author      struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	}
	Github string `yaml:"github"`
}

func NewConfig() Config {
	conf := Config{}
	conf.Author.Name = "Author Name"
	conf.Author.Email = "author@email.com"
	conf.Description = "Project Description"
	conf.ProjectName = "Project Name"
	return conf
}

var config = NewConfig()
var Metadata = map[string]string{}

var LinkRegex = `\[.*\]\(.*\)`
var ImageRegex = `!\[.*\]\(.*\)`
var BoldRegex = `\*\*.*\*\*`
var ItalicRegex = `\*.*\*`
var InlineCodeRegex = "`.*`"

var regexForInlineCode = regexp.MustCompile("^.*(`.+`)+.*(`.+`)*$")
var regexForLinks = regexp.MustCompile(`\[.*\]\(.*\)`)
var regexForImages = regexp.MustCompile(`!\[.*\]\(.*\)`)
var regexForBold = regexp.MustCompile(`\*\*.*\*\*`)
var regexForItalic = regexp.MustCompile(`\*.*\*`)

var regexForFullLineTags = regexp.MustCompile(fmt.Sprintf(`^((%s)|(%s)|(%s)|(%s)|(%s)).{0,1}$`, LinkRegex, ImageRegex, BoldRegex, ItalicRegex, InlineCodeRegex))

func LoadConfig(configFileName string) error {
	if configFileName == "" {
		fmt.Println("No config file provided, using default config")
		return nil
	}
	file, err := os.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	err = yaml.Unmarshal(file, &config)
	Metadata = map[string]string{
		"authorName":  config.Author.Name,
		"authorEmail": config.Author.Email,
		"description": config.Description,
		"projectName": config.ProjectName,
		"version":     config.Version,
		"github":      config.Github,
		"favicon":     config.Favicon,
		"pageTitle":   "Page Title",
	}
	if err != nil {
		return fmt.Errorf("error unmarshalling yaml: %v", err)
	}
	return nil
}

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
	for k, v := range Metadata {
		(*metadataValues)[k] = v
	}
	if fileLines[0] == "---" {
		for j = 1; j < len(fileLines); j++ {
			if (fileLines)[j] == "---" {
				break
			}
			prefix := strings.Split(fileLines[j], ": ")[0]
			(*metadataValues)[prefix] = fileLines[j][len(prefix)+2:]
		}
	}
	return j
}

func HandleLists(fileLines *[]string, i int, k string) int {
	lines := *fileLines
	if i == 0 || (!strings.HasPrefix(lines[i-1], "<ul>") && !strings.HasPrefix(lines[i-1], "<li>")) {
		lines[i] = "<ul><li>" + lines[i][2:] + "</li>"
	} else {
		lines[i] = "<li>" + lines[i][2:] + "</li>"
	}
	if i == len(lines)-1 || (!strings.HasPrefix(lines[i+1], "- ") && !strings.HasPrefix(lines[i+1], "* ")) {
		lines[i] = lines[i] + "</ul>"
	}
	return i
}

func HandleInlineCode(lineContent string) string {
	tickCount := strings.Count(lineContent, "`")
	for i := 0; i < tickCount/2; i++ {
		lineContent = html.EscapeString(lineContent)
		lineContent = strings.Replace(lineContent, "`", "<code>", 1)
		lineContent = strings.Replace(lineContent, "`", "</code>", 1)
		lineContent = strings.Replace(lineContent, "<code></code>", "``", -1)
		lineContent = strings.Replace(lineContent, "</code><code>", "``", -1)
	}
	return lineContent
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
	codeBlock = html.EscapeString(codeBlock)
	(*fileLines)[i] = "<pre>" + codeTag + codeBlock + "</code></pre>"
	(*fileLines)[i] = strings.TrimSpace((*fileLines)[i])
	if j != len((*fileLines))-1 {
		(*fileLines) = append((*fileLines)[:i+1], (*fileLines)[j+1:]...)
	} else {
		(*fileLines) = (*fileLines)[:i+1]
	}
}

func HandleLinks(lineContent string) string {
	regexForLinks := regexp.MustCompile(`\[.*\]\(.*\)`)
	linkPart := regexForLinks.FindAllString(lineContent, -1)
	for _, linkInMD := range linkPart {
		splitLink := strings.Split(linkInMD, "](")
		linkText := splitLink[0][1:]
		link := splitLink[1][:len(splitLink[1])-1]
		lineContent = strings.Replace(lineContent, linkInMD, "<a href=\""+link+"\" target=\"_blank\">"+linkText+"</a>", 1)
	}
	return lineContent
}

func HandleImages(lineContent string) string {
	images := regexForImages.FindAllString(lineContent, -1)
	for _, imageInMD := range images {
		imageInMD = strings.Trim(imageInMD, " ")
		splitLink := strings.Split(imageInMD, "](")
		endOfLink := strings.Index(splitLink[1], ")")
		splitLink[1] = splitLink[1][:endOfLink+1]
		altText := splitLink[0][2:]
		link := splitLink[1][:len(splitLink[1])-1]
		lineContent = strings.Replace(lineContent, imageInMD, "<img src=\""+link+"\" alt=\""+altText+"\">", 1)
	}
	return lineContent
}

func HandleBold(lineContent string) string {
	matches := regexForBold.FindAllString(lineContent, -1)
	for _, match := range matches {
		matchInHTML := match
		strongCount := strings.Count(match, "**")
		for i := 0; i < strongCount/2; i++ {
			matchInHTML = strings.Replace(matchInHTML, "**", "<strong>", 1)
			matchInHTML = strings.Replace(matchInHTML, "**", "</strong>", 1)
			matchInHTML = strings.Replace(matchInHTML, "<strong></strong>", "**", -1)
			matchInHTML = strings.Replace(matchInHTML, "</strong><strong>", "**", -1)
		}
		lineContent = strings.Replace(lineContent, match, matchInHTML, 1)
	}
	return lineContent
}

func HandleItalic(lineContent string) string {
	matches := regexForItalic.FindAllString(lineContent, -1)
	for _, match := range matches {
		matchInHTML := match
		italicCount := strings.Count(matchInHTML, "*")
		for i := 0; i < italicCount/2; i++ {
			matchInHTML = strings.Replace(matchInHTML, "*", "<em>", 1)
			matchInHTML = strings.Replace(matchInHTML, "*", "</em>", 1)
			matchInHTML = strings.Replace(matchInHTML, "<em></em>", "*", -1)
			matchInHTML = strings.Replace(matchInHTML, "</em><em>", "*", -1)
		}
		lineContent = strings.Replace(lineContent, match, matchInHTML, 1)
	}
	return lineContent
}

func MatchAndReplace(lineContent string) string {
	if regexForInlineCode.MatchString(lineContent) {
		lineContent = HandleInlineCode(lineContent)
	}
	if regexForBold.MatchString(lineContent) {
		lineContent = HandleBold(lineContent)
	}
	if regexForItalic.MatchString(lineContent) {
		lineContent = HandleItalic(lineContent)
	}
	if regexForLinks.MatchString(lineContent) && !regexForImages.MatchString(lineContent) {
		lineContentInHTML := HandleLinks(lineContent)
		if regexForFullLineTags.MatchString(lineContent) {
			return "<p>" + lineContentInHTML + "</p>"
		}
		return lineContentInHTML
	}
	if regexForImages.MatchString(lineContent) {
		lineContentInHTML := HandleImages(lineContent)
		if regexForFullLineTags.MatchString(lineContent) {
			return "<p>" + lineContentInHTML + "</p>"
		}
		return lineContentInHTML
	}

	return lineContent
}
