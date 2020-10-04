package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// export BEARER_TOKEN=<bearer token>

func auth() string {
	return os.Getenv("BEARER_TOKEN")
}

func createURL() string {
	return "https://api.twitter.com/2/tweets/search/stream"
}

func createRulesURL() string {
	return "https://api.twitter.com/2/tweets/search/stream/rules"
}

func createHeaders(bearerToken string) (string, string) {
	return "Authorization", "Bearer " + bearerToken
}

// getRules gets the set filter rules
func getRules(bearerToken string) []byte {
	req, _ := http.NewRequest("GET", createRulesURL(), nil)
	req.Header.Set(createHeaders(bearerToken))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("getRules():", string(bytes))
	return bytes
}

// deleteRules deletes the filter rule given as an argument
func deleteAllRules(bearerToken string, rules []byte) {
	if rules == nil {
		return
	}

	var rs struct {
		Datas []struct {
			ID string `json:"id"`
			//Value string `json:"value"`
			//Tag   string `json:"tag"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rules, &rs); err != nil {
		log.Fatal(err)
		return
	}
	ids := []string{}
	for _, data := range rs.Datas {
		ids = append(ids, data.ID)
	}

	tt := map[string]map[string][]string{
		"delete": {
			"ids": ids,
		},
	}
	reqBody, _ := json.Marshal(tt)
	req, _ := http.NewRequest("POST", createRulesURL(), bytes.NewBuffer(reqBody))
	req.Header.Set(createHeaders(bearerToken))
	req.Header.Set("Content-type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}

// setRules sets the filter rule
func setRules(bearerToken string) {
	// You can adjust the rules if needed
	t := []map[string]string{
		{"value": "dog has:images", "tag": "dog pictures"},
		{"value": "cat has:images -grumpy", "tag": "cat pictures"},
	}
	tt := map[string][]map[string]string{}
	tt["add"] = t
	reqBody, _ := json.Marshal(tt)

	req, _ := http.NewRequest("POST", createRulesURL(), bytes.NewBuffer(reqBody))
	req.Header.Set(createHeaders(bearerToken))
	req.Header.Set("Content-type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("setRules():", string(bytes))
}

func getStream(bearerToken string) {
	req, _ := http.NewRequest("GET", createURL(), nil)
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
	rules := getRules(bearerToken)
	deleteAllRules(bearerToken, rules)
	setRules(bearerToken)
	getStream(bearerToken)
}
