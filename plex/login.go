package plex

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Login(username, password string) (token string, error error) {
	client := &http.Client{}

	form := url.Values{}
	form.Add("user[login]", username)
	form.Add("user[password]", password)

	req, err := http.NewRequest("POST", "https://plex.tv/users/sign_in.json",
		strings.NewReader(form.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-plex-Client-Identifier", "MyWebHookServerV1")
	req.Header.Set("X-plex-Product", "WebHookServerV1")
	req.Header.Set("X-plex-Version", "V1")

	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return "", clientErr
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var v map[string]interface{}
	decodeErr := decoder.Decode(&v)
	if decodeErr != nil {
		return "", decodeErr
	}

	// Extract the authentication token from the returned JSON
	auth := v["user"].(map[string]interface{})
	aToken := auth["authToken"].(string)

	log.Printf("plex auth token: [%s]", aToken)

	return aToken, nil
}
