package kafka

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type FlagGenerator struct {
	OutputYamlFile   string
	TemplateYamlFile string
	NumberOfTopics   string
}

var flagGenerator = FlagGenerator{
	OutputYamlFile:   "output_yaml_file",
	TemplateYamlFile: "template_yaml_file",
	NumberOfTopics:   "num_of_topics",
}

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

func appendToYAML(templateFileName string, filename string, data SampleTopic) error {
	// Read the existing content of the YAML file
	existingContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("\n\n[ERROR]: Output file needs to be exist to append\n")
		return err
	}

	templateFile, err := os.ReadFile(templateFileName)
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

	template_yaml_file, _ := cmd.Flags().GetString(flagGenerator.TemplateYamlFile)

	output_yaml_file, _ := cmd.Flags().GetString(flagGenerator.OutputYamlFile)

	num_of_topics, _ := cmd.Flags().GetString(flagGenerator.NumberOfTopics)

	num_of_topics_to_int, err := strconv.Atoi(num_of_topics)

	if err != nil {
		fmt.Println("Error during conversion")
		return err
	}

	outputYamlFilename := string(output_yaml_file)

	for i := 0; i < num_of_topics_to_int; i++ {
		rand.Seed(time.Now().UnixNano())
		randomString := generateRandomString(10) //

		topic_data := SampleTopic{Topic_name: randomString, Replication_factor: 1, Partition_size: 1}

		err := appendToYAML(template_yaml_file, outputYamlFilename, topic_data)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		autoFormatYAML(topic_data)
	}

	fmt.Println("Template appended to", outputYamlFilename)

	return nil
}

// NOTE command definitions
var generateSampleTopicYaml = &cobra.Command{
	Use:   "generate-sample-topic-config", // camelCase
	Short: "",
	RunE:  generateTemplate,
}

func init() {

	generateSampleTopicYaml.Flags().String(flagGenerator.OutputYamlFile, "", "output yaml file")

	generateSampleTopicYaml.Flags().String(flagGenerator.TemplateYamlFile, "./template/kafka/topic.yaml", "template yaml file")

	generateSampleTopicYaml.Flags().String(flagGenerator.NumberOfTopics, "", "yaml file path")

	Topic.AddCommand(generateSampleTopicYaml)
}
