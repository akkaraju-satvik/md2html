package cmd

import (
	"fmt"
	"log"
	"md2htm/lib"
	"os"
	"path"
	"strings"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "md2htm",
	Short: "Converts markdown files to html",
	Run:   convert,
}

func convert(cmd *cobra.Command, args []string) {
	inputFileName := cmd.Flag("file").Value.String()
	fileNameWithoutExtension := strings.Split(inputFileName, ".")[0]
	inputFileExtension := strings.Split(inputFileName, ".")[1]
	if inputFileExtension != "md" && inputFileExtension != "markdown" {
		log.Fatal("File type not supported")
	}
	projectConfigFileName := cmd.Flag("config-file").Value.String()
	if err := lib.LoadConfigAndHandleCustomData(projectConfigFileName); err != nil {
		log.Fatal(err)
	}
	output := cmd.Flag("output").Value.String()
	if output == "" {
		output = lib.Metadata["outputDir"].(string) + "/" + fileNameWithoutExtension + ".html"
	}
	templateFileName := cmd.Flag("template-file").Value.String()
	var templateFile string
	if templateFileName == "" {
		fmt.Println("No template file provided, using default template file...")
		templateFile = lib.HtmlTemplate
	} else {
		templateFileData, err := os.ReadFile(templateFileName)
		if err != nil {
			log.Fatal(err)
		}
		templateFile = string(templateFileData)
	}
	convertedFileData, err := compile(inputFileName, templateFile)
	if err != nil {
		log.Fatal(err)
	}
	err = lib.CopyAssets(output, lib.Metadata["assetsDir"].(string))
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(output, []byte(convertedFileData), 0644)
	if err != nil {
		if os.IsNotExist(err) {
			dir := path.Dir(output)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatal(err)
			}
			err = os.WriteFile(output, []byte(convertedFileData), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File converted successfully")
}

func compile(inputFileName string, templateFile string) (string, error) {
	file, err := os.ReadFile(inputFileName)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	fileContents := string(file)
	fileLines := strings.Split(fileContents, "\n")
	if fileLines[len(fileLines)-1] == "" {
		fileLines = fileLines[:len(fileLines)-1]
	}
	var metadataValues = make(map[string]interface{})
	j := lib.HandleMetadata(fileLines, &metadataValues)
	for i := 0; i < len(fileLines); i++ {
		tagFound := false
		fileLines[i] = strings.TrimSpace(fileLines[i])
		if fileLines[i] == "" {
			tagFound = true
		}
		if i == 0 && fileLines[i] == "---" {
			fileLines = fileLines[j+1:]
			tagFound = true
		}
		if fileLines[i] == "---" {
			fileLines[i] = "<hr/>"
		}
		prefix := strings.Split(fileLines[i], " ")[0]
		fileLines[i] = lib.MatchAndReplace(fileLines[i])
		if lib.Tags[prefix] != "" || strings.Contains(fileLines[i], "```") {
			if prefix == "-" || prefix == "*" {
				lib.HandleLists(&fileLines, i, prefix)
			} else if strings.Contains(fileLines[i], "```") {
				lib.HandleCodeBlocks(&fileLines, i)
			} else {
				if i > 0 && (strings.HasPrefix(fileLines[i-1], "- ") || strings.HasPrefix(fileLines[i-1], "* ")) {
					fileLines[i-1] = fileLines[i-1] + "</ul>"
				}
				fileLines[i] = lib.ConvertToHTMLTags(prefix, fileLines[i])
			}
			tagFound = true
		}
		if !tagFound {
			lib.HandleParagraphs(&fileLines, i)
			tagFound = true
		}
	}
	htmlFormatOfFile := strings.Join(fileLines, "\n")
	templateFile = strings.Replace(templateFile, "$data", htmlFormatOfFile, 1)
	for k, v := range metadataValues {
		templateFile = strings.Replace(templateFile, "$"+k, v.(string), -1)
	}
	return templateFile, nil
}

func Execute() {
	rootCmd.Flags().StringP("file", "f", "", "The file to convert")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringP("output", "o", "", "The output file")
	rootCmd.Flags().StringP("config-file", "c", "", "Project Configuration file")
	rootCmd.Flags().StringP("template-file", "t", "", "Template file to use for conversion")
	rootCmd.Execute()
}
