package lib

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

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

func newConfig() Config {
	conf := Config{}
	conf.Author.Name = "Author Name"
	conf.Author.Email = "author@email.com"
	conf.Description = "Project Description"
	conf.ProjectName = "Project Name"
	conf.Version = "1.0.0"
	conf.Favicon = "favicon.ico"
	conf.Github = "https://github.com"
	return conf
}

var conf = newConfig()

func LoadConfigAndHandleCustomData(projectConfigFileName string, customDataFile string) error {
	if err := loadConfig(projectConfigFileName); err != nil {
		return err
	}
	if customDataFile != "" {
		if err := handleCustomData(customDataFile); err != nil {
			return err
		}
	}
	return nil
}

func loadConfig(configFileName string) error {
	if configFileName == "" {
		fmt.Println("No config file provided, using default config")
		return nil
	}
	file, err := os.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	err = yaml.Unmarshal(file, &conf)
	Metadata = map[string]string{
		"authorName":  conf.Author.Name,
		"authorEmail": conf.Author.Email,
		"description": conf.Description,
		"projectName": conf.ProjectName,
		"version":     conf.Version,
		"github":      conf.Github,
		"favicon":     conf.Favicon,
		"pageTitle":   "Page Title",
	}
	if err != nil {
		return fmt.Errorf("error unmarshalling yaml: %v", err)
	}
	return nil
}

func handleCustomData(fileName string) error {
	if !strings.HasSuffix(fileName, ".yaml") {
		return fmt.Errorf("file is not a yaml file")
	}
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	var customData map[string]string
	err = yaml.Unmarshal(fileContents, &customData)
	if err != nil {
		return fmt.Errorf("error unmarshalling yaml: %v", err)
	}
	for k, v := range customData {
		Metadata[k] = v
	}
	return nil
}
