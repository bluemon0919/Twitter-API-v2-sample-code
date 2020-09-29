package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

// export BEARER_TOKEN=<bearer token>

func auth() string {
	return os.Getenv("BEARER_TOKEN")
}

func createURL() string {
	return "https://api.twitter.com/2/tweets/sample/stream"
}

func createHeaders(bearerToken string) (string, string) {
	return "Authorization", "Bearer " + bearerToken
}

func connectToEndpoint(url string, bearerToken string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set(createHeaders(bearerToken))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(string(scanner.Bytes()))
	}
	return
}

func main() {
	bearerToken := auth()
	url := createURL()
	connectToEndpoint(url, bearerToken)
}
