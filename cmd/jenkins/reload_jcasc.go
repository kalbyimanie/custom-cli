package jenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func checkPrefixSlash(url string) (prefix string) {
	tmp := []string{}
	tmp = append(tmp, url)
	prefix = "/"

	for _, list := range tmp {
		lastIndex := strings.LastIndex(list, prefix)

		if list[lastIndex:] != prefix {
			return ""
		} else {
			return prefix
		}
	}

	return ""

}

type CrumbIssuer struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

func reloadJcasc(cmd *cobra.Command, args []string) error {
	// NOTE: format -> http://your-jenkins-url/reload-configuration-as-code

	// NOTE: Required flags
	jenkinsUrl, _ := cmd.Flags().GetString("jenkins-url")

	// NOTE: Required flags

	username := "admin"
	password := "admin"

	jBaseUrl, _ := cmd.Flags().GetString("jenkins-url")

	checkPrefix := checkPrefixSlash(jBaseUrl)

	if checkPrefix == "" {
		jBaseUrl = jBaseUrl + "/reload-configuration-as-code"
	} else {
		jBaseUrl = jBaseUrl + "reload-configuration-as-code"
	}

	fmt.Println(jBaseUrl)

	// NOTE: create httpClient
	client := &http.Client{}

	// NOTE: create timeout for httpClient
	context, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	// NOTE: defer cancel upon timeout
	defer cancel()

	// REVIEW: crumbIssuer section

	// NOTE: Fetch CSRF crumb
	crumbURL := jenkinsUrl + "/crumbIssuer/api/json"
	req, err := http.NewRequestWithContext(context, "GET", crumbURL, nil)
	if err != nil {
		fmt.Println("Error creating crumb request:", err)
		return err
	}

	req.SetBasicAuth(username, password)

	crumbResp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching crumb:", err)
		return err
	}
	defer crumbResp.Body.Close()

	// NOTE: Read the crumb response

	responseBody, err := ioutil.ReadAll(crumbResp.Body)

	if err != nil {
		fmt.Println("Error reading crumb response:", err)
		return err
	}

	var crumb CrumbIssuer
	err = json.Unmarshal(responseBody, &crumb)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	// NOTE: Print the crumb and crumbRequestField values
	fmt.Println(crumb.Crumb)
	fmt.Println(crumb.CrumbRequestField)

	// REVIEW: client making POST request  section
	// NOTE provide a client http post request with the given context
	req, err = http.NewRequestWithContext(context, "POST", jBaseUrl, nil)

	// NOTE: Set the Basic Authentication header upon the request
	req.SetBasicAuth(username, password)

	//// NOTE: set crumb data upon post request
	// //req.Header.Set(crumb.CrumbRequestField, crumb.Crumb)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return err
	}

	// NOTE: client making http post request

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return err
	}

	// NOTE: defer close upon client request once everything is done

	defer resp.Body.Close()

	// NOTE reading server response
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("%s\n", err)
	}

	fmt.Println(string(body))

	return nil

}

// NOTE command variable needs to be exported
var Jcasc = &cobra.Command{
	Use:   "jcasc", // camelCase
	Short: "",
	RunE:  reloadJcasc,
}

func init() {
	//NOTE set subcommands
	Jcasc.Flags().String("jenkins-url", "http://localhost/jenkins", "jenkins URL")
	Jcasc.MarkFlagRequired("jenkins-url")

	JenkinsCmd.AddCommand(Jcasc)
}
