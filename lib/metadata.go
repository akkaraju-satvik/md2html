package lib

import (
	"errors"
	"os"
	"path"
	"strings"
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
			background-color: #e0e0e0;
			font-family: monospace;
			border-radius: 5px;
		}
		code {
			font-family: monospace;
			background-color: #e0e0e0;
			padding: 0.2em;
		}
		pre code span {
			font-family: monospace;
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

var Metadata = map[string]interface{}{}

func HandleMetadata(fileLines []string, metadataValues *map[string]interface{}) int {
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

func CopyAssets(outputFileName string, assetsDir string) error {
	if assetsDir == "" {
		return nil
	}
	files, err := os.ReadDir(assetsDir)
	assetsDir = path.Clean(assetsDir)
	outputDir := path.Dir(outputFileName)
	os.RemoveAll(outputDir + "/" + assetsDir)
	if err != nil {
		return errors.New("error reading assets directory: ")
	}
	_ = os.MkdirAll(outputDir+"/"+assetsDir, 0755)
	done := make(chan bool)
	for _, file := range files {
		go func(file os.DirEntry) error {
			if file.IsDir() {
				CopyAssets(outputFileName, assetsDir+"/"+file.Name())
			} else {
				outputAssetFile := outputDir + "/" + assetsDir + "/" + file.Name()
				fileContents, err := os.ReadFile(assetsDir + "/" + file.Name())
				if err != nil {
					return errors.New("error reading file: ")
				}
				err = os.WriteFile(outputAssetFile, fileContents, 0644)
				if err != nil {
					return errors.New("error writing file: ")
				}
			}
			done <- true
			return nil
		}(file)
	}
	for range files {
		<-done
	}
	return nil
}
