package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gojektech/heimdall/v6/httpclient"
)

type caseData struct {
	dailyconfirmed string
	dailydeceased  string
	dailyrecovered string
	date           string
	totalconfirmed string
	totaldeceased  string
	totalrecovered string
}
type Stats struct {
	cases     []caseData `json:cases_time_series`
	statewise []interface{}
	tested    []interface{}
}

func main() {

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	fmt.Println("Hello world")
	fmt.Println("Triggered GH action")
	timeout := 1000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	// Use the clients GET method to create and execute the request
	res, err := client.Get("https://api.covid19india.org/data.json", nil)
	if err != nil {
		panic(err)
	}

	// Heimdall returns the standard *http.Response object
	body, err := ioutil.ReadAll(res.Body)

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	fmt.Println(string(body)[:100])
	if err != nil {
		fmt.Println(err)
	}

	tempdata := data["cases_time_series"].([]interface{})
	currentData := tempdata[len(tempdata)-1].(map[string]interface{})
	cd := make(map[string]string)
	for k, v := range currentData {
		cd[k] = v.(string)
	}
	fmt.Printf("\n\n %+v \n", cd["totalconfirmed"])

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's ACCOUNT:\n%+v\n", user)

	// Home Timeline
	tweets, resp, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 20,
	})

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update("just setting up my twttr", nil)

}
