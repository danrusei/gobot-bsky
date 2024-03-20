package gobotbsky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ResolveResponse struct {
	Did string `json:"did"`
}

type ApiKeyResponse struct {
	AccessJwt string `json:"accessJwt"`
}

type PostResponse struct {
	Feed string `json:"feed"`
}

func ResolveDID(handle string) (string, error) {
	resolveURL := "https://bsky.social/xrpc/com.atproto.identity.resolveHandle"
	params := map[string]string{"handle": handle}
	resp, err := http.Get(resolveURL + "?" + encodeParams(params))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resolveResp ResolveResponse
	err = json.Unmarshal(body, &resolveResp)
	if err != nil {
		return "", err
	}

	return resolveResp.Did, nil
}

func GetApiKey(identifier, password string) (string, error) {
	apiKeyURL := "https://bsky.social/xrpc/com.atproto.server.createSession"
	data := map[string]string{"identifier": identifier, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(apiKeyURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiKeyResp ApiKeyResponse
	err = json.Unmarshal(body, &apiKeyResp)
	if err != nil {
		return "", err
	}

	return apiKeyResp.AccessJwt, nil
}

func GetUserFeed(actor string, limit int, apiKey string) (string, error) {
	feedURL := "https://bsky.social/xrpc/app.bsky.feed.getAuthorFeed"
	params := map[string]interface{}{"actor": actor, "limit": limit}
	resp, err := httpGetWithAuth(feedURL, params, apiKey)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var postResp PostResponse
	err = json.Unmarshal(body, &postResp)
	if err != nil {
		return "", err
	}

	return postResp.Feed, nil
}

func PostToFeed(did, text string, apiKey string) (string, error) {
	postFeedURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
	record := map[string]interface{}{
		"collection": "app.bsky.feed.post",
		"repo":       did,
		"record": map[string]interface{}{
			"text":      text,
			"createdAt": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			"$type":     "app.bsky.feed.post",
		},
	}
	jsonData, err := json.Marshal(record)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(postFeedURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var postResp PostResponse
	err = json.Unmarshal(body, &postResp)
	if err != nil {
		return "", err
	}

	return postResp.Feed, nil
}

func httpGetWithAuth(url string, params map[string]interface{}, apiKey string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+apiKey)

	return client.Do(req)
}

func encodeParams(params map[string]string) string {
	var buf bytes.Buffer
	for key, value := range params {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(value)
	}
	return buf.String()
}
