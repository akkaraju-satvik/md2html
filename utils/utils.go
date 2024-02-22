package utils

import (
	"fmt"
	"os"
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
	"`":   "code",
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
  <title>$pageTitle | $projectName</title>
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

type Config struct {
	ProjectName string `yaml:"projectName"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
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

func LoadConfig(configFileName string) (Config, error) {
	if configFileName == "" {
		fmt.Println("No config file provided, using default config")
		return config, nil
	}
	file, err := os.ReadFile(configFileName)
	if err != nil {
		return config, fmt.Errorf("error reading file: %v", err)
	}
	err = yaml.Unmarshal(file, &config)
	Metadata = map[string]string{
		"authorName":  config.Author.Name,
		"authorEmail": config.Author.Email,
		"description": config.Description,
		"projectName": config.ProjectName,
		"version":     config.Version,
		"github":      config.Github,
		"pageTitle":   "Page Title",
	}
	if err != nil {
		return config, fmt.Errorf("error unmarshalling yaml: %v", err)
	}
	return config, nil
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
	if fileLines[0] == "---" {
		for k, v := range Metadata {
			(*metadataValues)[k] = v
		}
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
