package utils

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
