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

	// Facets Section
	// =======================================
	// Facet_type coulf be Facet_Link, Facet_Mention or Facet_Tag
	// based on the selected type it expect the second argument to be URI, DID, or TAG
	// the last function argument is the text, part of the original text that is modifiend in Richtext

	post1, err := gobot.NewPostBuilder("Hello to Bluesky, the coolest open social network").WithFacet(gobot.Facet_Link, "https://docs.bsky.app/", "Bluesky").WithFacet(gobot.Facet_Tag, "bsky", "open social").Build()
	if err != nil {
		fmt.Printf("Got error: %v", err)
	}

	cid1, uri1, err := agent.PostToFeed(ctx, post1)
	if err != nil {
		fmt.Printf("Got error: %v", err)
	} else {
		fmt.Printf("Succes: Cid = %v , Uri = %v", cid1, uri1)
	}

	// Embed Links section
	// =======================================

	u, err := url.Parse("https://go.dev/")
	if err != nil {
		log.Fatalf("Parse error, %v", err)
	}
	post2, err := gobot.NewPostBuilder("Hello to Go on Bluesky").WithExternalLink("Go Programming Language", *u, "Build simple, secure, scalable systems with Go").Build()
	if err != nil {
		fmt.Printf("Got error: %v", err)
	}

	cid2, uri2, err := agent.PostToFeed(ctx, post2)
	if err != nil {
		fmt.Printf("Got error: %v", err)
	} else {
		fmt.Printf("Succes: Cid = %v , Uri = %v", cid2, uri2)
	}

	// Embed Images section
	// =======================================
	images := []gobot.Image{}

	url1, err := url.Parse("https://www.freecodecamp.org/news/content/images/2021/10/golang.png")
	if err != nil {
		log.Fatalf("Parse error, %v", err)
	}
	images = append(images, gobot.Image{
		Title: "Golang",
		Uri:   *url1,
	})

	blobs, err := agent.UploadImages(ctx, images...)
	if err != nil {
		log.Fatalf("Parse error, %v", err)
	}

	post3, err := gobot.NewPostBuilder("Gobot-bsky - a simple golang lib to write Bluesky bots").WithImages(blobs, images).Build()
	if err != nil {
		fmt.Printf("Got error: %v", err)
	}

	cid3, uri3, err := agent.PostToFeed(ctx, post3)
	if err != nil {
		fmt.Printf("Got error: %v", err)
	} else {
		fmt.Printf("Succes: Cid = %v , Uri = %v", cid3, uri3)
	}

}
