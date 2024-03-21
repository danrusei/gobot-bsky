package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"

	gobot "github.com/danrusei/gobot-bsky"
)

func main() {

	godotenv.Load()
	handle := os.Getenv("HANDLE")
	apikey := os.Getenv("APIKEY")
	server := "https://bsky.social"

	ctx := context.Background()

	agent := gobot.NewAgent(ctx, server, handle, apikey)
	agent.Connect(ctx)

	u, err := url.Parse("https://go.dev/")
	if err != nil {
		log.Fatalf("Parse error, %v", err)
	}
	post := gobot.NewPostBuilder("Hello to Bluesky").WithExternalLink("Gopher", *u, "Build simple, secure, scalable systems with Go").Build()

	cid, uri, err := agent.PostToFeed(ctx, post)
	if err != nil {
		fmt.Printf("Got error: %v", err)
	} else {
		fmt.Printf("Succes: Cid = %v , Uri = %v", cid, uri)
	}

}
