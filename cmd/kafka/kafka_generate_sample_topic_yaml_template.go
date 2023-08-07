package kafka

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type SampleTopic struct {
	Topic_name         string
	Replication_factor int
	Partition_size     int
}

func autoFormatYAML(data interface{}) (string, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func appendToYAML(filename string, data SampleTopic) error {
	// Read the existing content of the YAML file
	existingContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// template file
	templateFile, err := os.ReadFile("./template/kafka/topic.yaml")
	if err != nil {
		fmt.Println("Error reading template:", err)
		return err
	}

	// Parse the template
	tmpl, err := template.New("yamlTemplate").Parse(string(templateFile))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

	// Prepare the template data
	var newData string
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}

	newData = buf.String()

	// Append the new template data to the existing content
	updatedContent := append(existingContent, []byte(newData)...)

	// Write the updated content back to the file
	err = ioutil.WriteFile(filename, updatedContent, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func generateTemplate(cmd *cobra.Command, args []string) error {
	yamlFilename := "/Users/kalbyimanie/Documents/personal/repositories/custom-cli/config/kafka/env/prod/hotel_topics.yaml"

	for i := 0; i < 900; i++ {
		rand.Seed(time.Now().UnixNano())
		randomString := generateRandomString(10) //

		topic_data := SampleTopic{Topic_name: randomString, Replication_factor: 1, Partition_size: 1}

		err := appendToYAML(yamlFilename, topic_data)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		autoFormatYAML(topic_data)
	}

	fmt.Println("Template appended to", yamlFilename)

	return nil
}

// NOTE command definitions
var generateSampleTopicYaml = &cobra.Command{
	Use:   "generate-sample-topic-config", // camelCase
	Short: "",
	RunE:  generateTemplate,
}

func init() {
	Topic.AddCommand(generateSampleTopicYaml)
}
