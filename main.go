package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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

	consumerKey := os.Getenv("consumer-key")
	consumerSecret := os.Getenv("consumer-secret")
	accessToken := os.Getenv("access-token")
	accessSecret := os.Getenv("access-secret")

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
	if err != nil {
		fmt.Println(err)
	}

	tempdata := data["cases_time_series"].([]interface{})
	currentData := tempdata[len(tempdata)-1].(map[string]interface{})
	cd := make(map[string]string)
	for k, v := range currentData {
		cd[k] = v.(string)
	}
	tweetString := fmt.Sprintf("Today's confirmed cases: %s\nToday's recoveries: %s\nDeceased today: %s\nTotal Confirmed cases: %s\nTotal recoveries: %s\nTotal casualties:%s", cd["dailyconfirmed"], cd["dailyrecovered"], cd["dailydeceased"], cd["totalconfirmed"], cd["totalrecovered"], cd["totaldeceased"])

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	twitterClient := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, _ := twitterClient.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's ACCOUNT:\n%+v\n", user)

	// Home Timeline
	// tweets, resp, err := twitterClient.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
	// 	Count: 20,
	// })

	// Send a Tweet
	_, _, err = twitterClient.Statuses.Update(tweetString, nil)
	fmt.Println(tweetString)
	fmt.Println(len(tweetString))
}
