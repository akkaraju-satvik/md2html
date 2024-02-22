package cmd

import (
	"fmt"
	"log"
	"md2htm/utils"
	"os"
	"strings"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "md2htm",
	Short: "Converts markdown files to html",
	Run: func(cmd *cobra.Command, args []string) {
		inputFileName := cmd.Flag("file").Value.String()
		fileNameWithoutExtension := strings.Split(inputFileName, ".")[0]
		inputFileExtension := strings.Split(inputFileName, ".")[1]
		if inputFileExtension != "md" && inputFileExtension != "markdown" {
			log.Fatal("File type not supported")
		}
		output := cmd.Flag("output").Value.String()
		if output == "" {
			output = fileNameWithoutExtension + ".html"
		}
		projectConfigFileName := cmd.Flag("config-file").Value.String()
		utils.LoadConfig(projectConfigFileName)
		convertedFileData, err := convert(inputFileName, output)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(output, []byte(convertedFileData), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("File converted successfully")
	},
}

func convert(inputFileName string, output string) (string, error) {
	file, err := os.ReadFile(inputFileName)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	fileContents := string(file)
	fileLines := strings.Split(fileContents, "\n")
	if fileLines[len(fileLines)-1] == "" {
		fileLines = fileLines[:len(fileLines)-1]
	}
	var metadataValues = make(map[string]string)
	for i := 0; i < len(fileLines); i++ {
		tagFound := false
		if fileLines[i] == "" {
			fileLines[i] = ""
			tagFound = true
		}
		if i == 0 && fileLines[i] == "---" {
			j := utils.HandleMetadata(fileLines, &metadataValues)
			fileLines = fileLines[j+1:]
			tagFound = true
		}
		for k := range utils.Tags {
			if strings.HasPrefix(fileLines[i], k) {
				if k == "- " || k == "* " {
					utils.HandleLists(&fileLines, i, k)
				} else if strings.Contains(fileLines[i], "```") {
					utils.HandleCodeBlocks(&fileLines, i)
				} else {
					if i > 0 && (strings.HasPrefix(fileLines[i-1], "- ") || strings.HasPrefix(fileLines[i-1], "* ")) {
						fileLines[i-1] = fileLines[i-1] + "</ul>"
					}
					fileLines[i] = utils.ConvertToHTMLTags(k, fileLines[i])
				}
				tagFound = true
			}
		}
		if !tagFound {
			fileLines[i] = "<p>" + fileLines[i] + "</p>"
			tagFound = true
		}
	}
	htmlFormatOfFile := strings.Join(fileLines, "\n")
	templateHtml := utils.HtmlTemplate
	if err != nil {
		panic(err)
	}
	for k, v := range metadataValues {
		templateHtml = strings.Replace(string(templateHtml), "$"+k, v, -1)
	}
	templateHtml = strings.Replace(string(templateHtml), "$data", htmlFormatOfFile, 1)
	return string(templateHtml), nil
}

func Execute() {
	rootCmd.Flags().StringP("file", "f", "", "The file to convert")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringP("output", "o", "", "The output file")
	rootCmd.Flags().StringP("config-file", "c", "", "Project Configuration file")
	rootCmd.Execute()
}
