package lib

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

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

func HandleParagraphs(fileLines *[]string, i int) {
	lastLineIndex := len((*fileLines)) - 1
	isLastLine := i == lastLineIndex
	if i-1 >= 0 && (*fileLines)[i-1] == "" {
		(*fileLines)[i] = "<p>" + (*fileLines)[i]
	}
	if i <= lastLineIndex {
		if !isLastLine && (*fileLines)[i+1] == "" {
			(*fileLines)[i] = (*fileLines)[i] + "</p>"
		} else if isLastLine && (*fileLines)[i] != "" {
			(*fileLines)[i] = (*fileLines)[i] + "</p>"
		}
	}
}
