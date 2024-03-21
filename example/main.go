package main

import (
	"context"
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

	post := gobot.NewPostBuilder()

	agent.PostToFeed(ctx, *post)

}
