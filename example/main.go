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

	// u, err := url.Parse("https://go.dev/")
	// if err != nil {
	// 	log.Fatalf("Parse error, %v", err)
	// }
	// post := gobot.NewPostBuilder("Hello to Bluesky").WithExternalLink("Gopher", *u, "Build simple, secure, scalable systems with Go").Build()

	images := []gobot.Image{}

	url1, err := url.Parse("https://www.freecodecamp.org/news/content/images/2021/10/golang.png")
	if err != nil {
		log.Fatalf("Parse error, %v", err)
	}
	images = append(images, gobot.Image{
		Title: "Golang",
		Uri:   *url1,
	})

	// url2, err := url.Parse("https://pkg.go.dev/static/shared/gopher/package-search-700x300.jpeg")
	// if err != nil {
	// 	log.Fatalf("Parse error, %v", err)
	// }
	// images = append(images, gobot.Image{
	// 	Title: "pkg.go.dev",
	// 	Uri:   *url2,
	// })

	blobs, err := agent.UploadImages(ctx, images...)
	if err != nil {
		log.Fatalf("Parse error, %v", err)
	}

	post := gobot.NewPostBuilder("Hello to Bluesky").WithImages(blobs, images).Build()

	cid, uri, err := agent.PostToFeed(ctx, post)
	if err != nil {
		fmt.Printf("Got error: %v", err)
	} else {
		fmt.Printf("Succes: Cid = %v , Uri = %v", cid, uri)
	}

}
