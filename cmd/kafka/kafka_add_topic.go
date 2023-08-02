package kafka

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NOTE flags definitions
type Flags struct {
	KafkaBrokerConfigYamlFile     string
	KafkaBrokerConfigYamlFilePath string
	KafkaTopicConfigYamlFile      string
	KafkaTopicConfigYamlFilePath  string
	Env                           string
}

// NOTE create flags from flags definitions
var flag = Flags{
	KafkaBrokerConfigYamlFile:     "broker-config",
	KafkaBrokerConfigYamlFilePath: "broker-config-path",
	KafkaTopicConfigYamlFile:      "topic-config",
	KafkaTopicConfigYamlFilePath:  "topic-config-path",
	Env:                           "env",
}

// NOTE kafka broker connection configuration
type BrokerConfig struct {
	Host string `mapstructure:"broker_host"`
	Port string `mapstructure:"broker_port"`
}

// NOTE: topic configuration
type TopicConfig struct {
	Name              string `mapstructure:"name"`
	ReplicationFactor int16  `mapstructure:"replication_factor"`
	PartitionSize     int32  `mapstructure:"partition_size"`
}

// NOTE root configurations
type Config struct {
	Topics       []TopicConfig  `mapstructure:"topics"`
	BrokerConfig []BrokerConfig `mapstructure:"brokers"`
}

func createTopic(cmd *cobra.Command, args []string) error {

	// NOTE env config
	env, _ := cmd.Flags().GetString(flag.Env)

	// path values
	defaultpathdev := "./config/kafka/env/dev/"
	defaultpathstage := "./config/kafka/env/stage/"
	defaultpathprod := "./config/kafka/env/prod/"
	bconfcustompath := ""
	tconfcustompath := ""

	// NOTE check path values
	bconfcheckpathvalue := cmd.Flags().Lookup(flag.KafkaBrokerConfigYamlFilePath).Value.String()

	tconfcheckpathvalue := cmd.Flags().Lookup(flag.KafkaTopicConfigYamlFilePath).Value.String()

	switch env {
	case "dev":
		if bconfcheckpathvalue == "" {
			bconfcustompath = defaultpathdev
		} else {
			// NOTE receive broker config file path
			bconfcustompath, _ = cmd.Flags().GetString(flag.KafkaBrokerConfigYamlFilePath)
		}

		if tconfcheckpathvalue == "" {
			tconfcustompath = defaultpathdev
		} else {
			// NOTE receive topic config file path
			tconfcustompath, _ = cmd.Flags().GetString(flag.KafkaTopicConfigYamlFilePath)
		}

	case "stage":
		if bconfcheckpathvalue == "" {
			bconfcustompath = defaultpathstage
		} else {
			// NOTE receive broker config file path
			bconfcustompath, _ = cmd.Flags().GetString(flag.KafkaBrokerConfigYamlFilePath)
		}

		if tconfcheckpathvalue == "" {
			tconfcustompath = defaultpathstage
		} else {
			// NOTE receive topic config file path
			tconfcustompath, _ = cmd.Flags().GetString(flag.KafkaTopicConfigYamlFilePath)
		}
	case "prod":
		if bconfcheckpathvalue == "" {
			bconfcustompath = defaultpathprod
		} else {
			// NOTE receive broker config file path
			bconfcustompath, _ = cmd.Flags().GetString(flag.KafkaBrokerConfigYamlFilePath)
		}

		if tconfcheckpathvalue == "" {
			tconfcustompath = defaultpathprod
		} else {
			// NOTE receive topic config file path
			tconfcustompath, _ = cmd.Flags().GetString(flag.KafkaTopicConfigYamlFilePath)
		}
	default:
		fmt.Fprintf(os.Stderr, "[ERROR] flag has to be 'dev' or 'stage' or 'prod'\n\n")
		cmd.Help()
		os.Exit(1)
	}

	// NOTE receive broker config flag
	bconf, _ := cmd.Flags().GetString(flag.KafkaBrokerConfigYamlFile)

	// NOTE receive topic config flag
	tconf, _ := cmd.Flags().GetString(flag.KafkaTopicConfigYamlFile)

	// NOTE prepare broker configuration yaml
	viper.AddConfigPath(bconfcustompath)
	viper.SetConfigName(bconf)
	viper.SetConfigType("yaml")

	// NOTE read broker configuration yaml
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading broker configuration yaml:", err)
		cmd.Help()
		os.Exit(1)
		return err
	}

	// NOTE Unmarshal the broker configuration yaml data into a struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Error unmarshaling config data:", err)
		return err
	}

	var broker_config_host string
	var broker_config_port string
	var broker_connection []string

	for _, broker_config := range config.BrokerConfig {
		broker_config_host = broker_config.Host
		broker_config_port = broker_config.Port

		propogate_broker := broker_config_host + ":" + broker_config_port

		broker_connection = append(broker_connection, propogate_broker)
	}

	// NOTE create broker connection
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V2_6_0_0 // Use the appropriate Kafka version
	kafkaConfig.ClientID = "kafka-topic-creator"

	admin, err := sarama.NewClusterAdmin(broker_connection, kafkaConfig)

	if err != nil {
		panic(fmt.Errorf("failed to connect to Kafka: %s\n\n", broker_connection))
	}
	defer admin.Close()

	// NOTE prepare topic configuration yaml
	viper.AddConfigPath(tconfcustompath)
	viper.SetConfigName(tconf)
	viper.SetConfigType("yaml")

	// NOTE read topic configuration yaml
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading broker configuration yaml:", err)
		cmd.Help()
		os.Exit(1)
		return err
	}

	// NOTE Unmarshal the topic configuration yaml data into a struct
	var configTopic Config
	if err := viper.Unmarshal(&configTopic); err != nil {
		fmt.Println("Error unmarshaling config data:", err)
		return err
	}

	// NOTE flag variable for go routine group counter
	var flag []string
	for _, topic_counter := range configTopic.Topics {
		flag = append(flag, topic_counter.Name)
	}

	var wg sync.WaitGroup

	// NOTE Create a channel to signal the main goroutine when processing is done
	done := make(chan struct{}, len(flag))

	// NOTE create topic
	for _, topic_config := range configTopic.Topics {
		// NOTE process counter and increment the counter by 1
		wg.Add(1)
		// NOTE start processing
		go processingTopic(topic_config.Name, topic_config.ReplicationFactor, topic_config.PartitionSize, admin, &wg, done)
	}
	// NOTE wait until each process finished
	wg.Wait()

	// NOTE close once all the process signal have been received
	close(done)

	// REVIEW Wait for the main goroutine to receive all done signals
	// for i := 0; i < len(flag); i++ {
	// 	<-done
	// }

	return nil
}

func processingTopic(t string, r int16, p int32, admin sarama.ClusterAdmin, wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()
	var err error

	err = admin.CreateTopic(t, &sarama.TopicDetail{
		NumPartitions:     p,
		ReplicationFactor: r,
	}, false)

	if err != nil {
		log.Printf("Topic '%s' has already been created.\n", err)
	} else {
		log.Printf("Topic '%s' created successfully.\n", t)
	}

	done <- struct{}{}
}

// NOTE command definitions
var createTopicCmd = &cobra.Command{
	Use:   "create", // camelCase
	Short: "",
	RunE:  createTopic,
}

func init() {
	// NOTE init flags
	createTopicCmd.Flags().String(flag.Env, "", "Env Name")

	createTopicCmd.Flags().String(flag.KafkaBrokerConfigYamlFilePath, "", "Broker Config YAML file path")
	createTopicCmd.Flags().String(flag.KafkaBrokerConfigYamlFile, "", "Broker Config YAML file")

	createTopicCmd.Flags().String(flag.KafkaTopicConfigYamlFilePath, "", "Topic Config YAML file path")
	createTopicCmd.Flags().String(flag.KafkaTopicConfigYamlFile, "", "Topic Config YAML file")

	// NOTE add command
	Topic.AddCommand(createTopicCmd)
}
