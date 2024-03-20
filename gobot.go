package gobotbsky

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

const defaultPDS = "https://bsky.social"

// const resolveURL = "https://bsky.social/xrpc/com.atproto.identity.resolveHandle"
// const sessionURL = "https://bsky.social/xrpc/com.atproto.server.createSession"
// const postFeedURL = "https://bsky.social/xrpc/com.atproto.repo.createRecord"

// Wrapper over the atproto xrpc transport
type BskyAgent struct {
	// xrpc transport, a wrapper around http server
	client *xrpc.Client
	handle string
	apikey string
}

// Creates new BlueSky Agent
func NewAgent(ctx context.Context, server string, handle string, apikey string) BskyAgent {

	if server == "" {
		server = defaultPDS
	}

	return BskyAgent{
		client: &xrpc.Client{
			Client: new(http.Client),
			Host:   server,
		},
		handle: handle,
		apikey: apikey,
	}

}

// Connect and Authenticate to the provided Personal Data Server, default is Bluesky PDS
// No need to refresh the access token if the bot script will be executed based on the cron job
func (c *BskyAgent) Connect(ctx context.Context) error {
	// Authenticate with the Bluesky server

	input_for_session := &atproto.ServerCreateSession_Input{
		Identifier: c.handle,
		Password:   c.apikey,
	}

	session, err := atproto.ServerCreateSession(ctx, c.client, input_for_session)

	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT: %v", err)
	}

	// Access Token is used to make authenticated requests
	// Refresh Token allows to generate a new Access Token
	c.client.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Handle:     session.Handle,
		Did:        session.Did,
	}

	return nil
}

func (c *BskyAgent) PostToFeed(ctx context.Context) error {

	post_input := &atproto.RepoCreateRecord_Input{
		// collection: The NSID of the record collection.
		Collection:
		// record: The record itself. Must contain a $type field.
		Record:
		// repo: The handle or DID of the repo (aka, current account).
		Repo: c.client.Auth.Did
		// rkey: The Record Key.
		Rkey: 
		// swapCommit: Compare and swap with the previous commit by CID.
		SwapCommit: 
		// validate: Can be set to 'false' to skip Lexicon schema validation of record data.
		Validate: true
	}
	
	
	response, err := atproto.RepoCreateRecord(ctx, c.client, post_input)

}

// type ResolveResponse struct {
// 	Did string `json:"did"`
// }

// type ApiKeyResponse struct {
// 	AccessJwt string `json:"accessJwt"`
// }

// type PostResponse struct {
// 	Feed string `json:"feed"`
// }

// func ResolveDID(handle string) (string, error) {
// 	//resolveURL := "https://bsky.social/xrpc/com.atproto.identity.resolveHandle"
// 	params := map[string]string{"handle": handle}
// 	resp, err := http.Get(resolveURL + "?" + encodeParams(params))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var resolveResp ResolveResponse
// 	err = json.Unmarshal(body, &resolveResp)
// 	if err != nil {
// 		return "", err
// 	}

// 	return resolveResp.Did, nil
// }

// func GetApiKey(identifier, password string) (string, error) {
// 	//apiKeyURL := "https://bsky.social/xrpc/com.atproto.server.createSession"
// 	data := map[string]string{"identifier": identifier, "password": password}
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return "", err
// 	}

// 	resp, err := http.Post(sessionURL, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var apiKeyResp ApiKeyResponse
// 	err = json.Unmarshal(body, &apiKeyResp)
// 	if err != nil {
// 		return "", err
// 	}

// 	return apiKeyResp.AccessJwt, nil
// }

// func GetUserFeed(actor string, limit int, apiKey string) (string, error) {
// 	feedURL := "https://bsky.social/xrpc/app.bsky.feed.getAuthorFeed"
// 	params := map[string]interface{}{"actor": actor, "limit": limit}
// 	resp, err := httpGetWithAuth(feedURL, params, apiKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var postResp PostResponse
// 	err = json.Unmarshal(body, &postResp)
// 	if err != nil {
// 		return "", err
// 	}

// 	return postResp.Feed, nil
// }

// func PostToFeed(did, text string, apiKey string) (string, error) {
// 	//postFeedURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
// 	record := map[string]interface{}{
// 		"collection": "app.bsky.feed.post",
// 		"repo":       did,
// 		"record": map[string]interface{}{
// 			"text":      text,
// 			"createdAt": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
// 			"$type":     "app.bsky.feed.post",
// 		},
// 	}
// 	jsonData, err := json.Marshal(record)
// 	if err != nil {
// 		return "", err
// 	}

// 	resp, err := http.Post(postFeedURL, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var postResp PostResponse
// 	err = json.Unmarshal(body, &postResp)
// 	if err != nil {
// 		return "", err
// 	}

// 	return postResp.Feed, nil
// }

// func httpGetWithAuth(url string, params map[string]interface{}, apiKey string) (*http.Response, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	q := req.URL.Query()
// 	for key, value := range params {
// 		q.Add(key, fmt.Sprintf("%v", value))
// 	}
// 	req.URL.RawQuery = q.Encode()

// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	return client.Do(req)
// }

// func encodeParams(params map[string]string) string {
// 	var buf bytes.Buffer
// 	for key, value := range params {
// 		if buf.Len() > 0 {
// 			buf.WriteByte('&')
// 		}
// 		buf.WriteString(key)
// 		buf.WriteByte('=')
// 		buf.WriteString(value)
// 	}
// 	return buf.String()
// }
