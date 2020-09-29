package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// export BEARER_TOKEN=<bearer token>

func auth() string {
	return os.Getenv("BEARER_TOKEN")
}

func createURL() string {
	// Specify the usernames that you want to lookup below
	// You can enter up to 100 comma-separated values.
	usernames := "usernames=TwitterDev,TwitterAPI,bluemon0919"
	// Tweet fields are adjustable.
	// See below for options:
	// https://developer.twitter.com/en/docs/twitter-api/tweets/search/api-reference/get-tweets-search-recent
	// user.fields
	userFields := "user.fields=description,created_at"
	url := fmt.Sprintf("https://api.twitter.com/2/users/by?%s&%s", usernames, userFields)
	return url
}

func createHeaders(bearerToken string) (string, string) {
	return "Authorization", "Bearer " + bearerToken
}

func connectToEndpoint(url string, bearerToken string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set(createHeaders(bearerToken))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes
}

func main() {
	bearerToken := auth()
	url := createURL()
	respBody := connectToEndpoint(url, bearerToken)
	fmt.Println(string(respBody))
}
