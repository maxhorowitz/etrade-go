package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
)

func main() {
	oauth()
}

func oauth() {

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	consKey := os.Getenv("ETRADE_SANDBOX_API_KEY")
	if consKey == "" {
		log.Fatal(errors.New("missing env variable"))
	}
	consSecret := os.Getenv("ETRADE_SANDBOX_API_SECRET")
	if consSecret == "" {
		log.Fatal(errors.New("missing env variable"))
	}

	config := oauth1.Config{
		ConsumerKey:    consKey,
		ConsumerSecret: consSecret,
		CallbackURL:    "oob",
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://api.etrade.com/oauth/request_token",
			AuthorizeURL:    "https://us.etrade.com/e/t/etws/authorize?key={}&token={}",
			AccessTokenURL:  "https://api.etrade.com/oauth/access_token",
		},
	}

	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		log.Fatal(err)
	}

	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("open this url in your browser: %s\n\n", authorizationURL.String())
	}

	fmt.Printf("Choose whether to grant the application access.\nPaste " +
		"the oauth_verifier parameter (excluding trailing #_=_) from the " +
		"address bar: ")
	var verifier string
	_, err = fmt.Scanf("%s", &verifier)
	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		log.Fatal("\n", err)
	}

	token := oauth1.NewToken(accessToken, accessSecret)
	fmt.Println(fmt.Sprint(token))

	httpClient := config.Client(oauth1.NoContext, token)

	path := "https://apisb.etrade.com/v1/market/quote/{AAPL}"
	res, err := httpClient.Get(path)
	if err != nil {
		fmt.Println("\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Printf("\nRaw Response Body:\n%v\n", string(body))

}
