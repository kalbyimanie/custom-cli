package network

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type flag struct {
	flagName string
	flagDesc string
}

var urlFlag = flag{
	flagName: "url",
	flagDesc: "example value (--url https://ipinfo.io/json)",
}

func geturl(cmd *cobra.Command, args []string) error {
	// NOTE flag url
	url, _ := cmd.Flags().GetString("url")

	client := &http.Client{}
	context, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// NOTE provide an http request method
	req, err := http.NewRequestWithContext(context, "GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return err
	}

	defer resp.Body.Close()

	// NOTE read response
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("%s\n", err)
	}

	fmt.Println(string(body))

	return nil

}

// command variable needs to be exported
var getUrl = &cobra.Command{
	Use:   "geturl --url <URL>", // camelCase
	Short: urlFlag.flagDesc,
	RunE:  geturl,
}

func init() {
	//NOTE set geturl subcommand
	getUrl.Flags().String(urlFlag.flagName, "", urlFlag.flagDesc)

	//NOTE add geturl subcommand to NetCmd
	NetworkCmd.AddCommand(getUrl)
}
