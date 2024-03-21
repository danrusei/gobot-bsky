package gobotbsky

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bluesky-social/indigo/api/atproto"
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"

	"github.com/bluesky-social/indigo/xrpc"
)

const defaultPDS = "https://bsky.social"

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

// Post to social app
func (c *BskyAgent) PostToFeed(ctx context.Context, post appbsky.FeedPost) (string, string, error) {

	post_input := &atproto.RepoCreateRecord_Input{
		// collection: The NSID of the record collection.
		Collection: "app.bsky.feed.post",
		// repo: The handle or DID of the repo (aka, current account).
		Repo: c.client.Auth.Did,
		// record: The record itself. Must contain a $type field.
		Record: &lexutil.LexiconTypeDecoder{Val: &post},
	}

	response, err := atproto.RepoCreateRecord(ctx, c.client, post_input)
	if err != nil {
		return "", "", fmt.Errorf("unable to post, %v", err)
	}

	return response.Cid, response.Uri, nil
}
