package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/garyburd/go-oauth/oauth"
)

// export CONSUMER_KEY=<your_consumer_key>
// export CONSUMER_SECRET=<your_consumer_secret>

func main() {
	oauthClient := &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  os.Getenv("CONSUMER_KEY"),
			Secret: os.Getenv("CONSUMER_SECRET"),
		},
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}

	scope := url.Values{"scope": {"read_public,write_public,read_private,write_private"}}
	// callbackがない場合は"oob"を設定する
	tempCredentials, err := oauthClient.RequestTemporaryCredentials(nil, "oob", scope)
	if err != nil {
		log.Fatal("RequestTemporaryCredentials:", err)
	}

	url := oauthClient.AuthorizationURL(tempCredentials, nil)
	fmt.Println("1. Go to ", url)
	fmt.Println("2. Authorize the application")
	fmt.Println("3. Enter verification code:")

	var code string
	fmt.Scanln(&code)
	fmt.Println("InputCode: ", code)

	credentials, _, err := oauthClient.RequestToken(nil, tempCredentials, code)
	if err != nil {
		log.Fatal("RequestToken:", err)
	}
	fmt.Println("Token: ", credentials.Token)
	fmt.Println("Secret: ", credentials.Secret)

	params := make(map[string][]string)
	params["ids"] = []string{"1278747501642657792"}
	params["tweet.fields"] = []string{"created_at"}
	resp, err := oauthClient.Get(nil, credentials,
		"https://api.twitter.com/2/tweets", params)
	if err != nil {
		log.Fatal("Get: ", err)
	}
	fmt.Println("Status: ", resp.Status)

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Read Err:", err)
	}
	fmt.Println(string(bytes))
}
