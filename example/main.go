package main

import (
	"fmt"

	gobot "github.com/danrusei/gobot-bsky"
)

func main() {
	handle := "totototo"
	appPassword := "myapppassword"
	actor := "tototot"
	limit := 4
	text := "Hello, world"

	did, err := gobot.ResolveDID(handle)
	if err != nil {
		fmt.Println("Error resolving DID:", err)
		return
	}

	apiKey, err := gobot.GetApiKey(did, appPassword)
	if err != nil {
		fmt.Println("Error getting API key:", err)
		return
	}

	userFeed, err := gobot.GetUserFeed(actor, limit, apiKey)
	if err != nil {
		fmt.Println("Error getting user feed:", err)
		return
	}

	fmt.Println("User feed:", userFeed)

	if _, err := gobot.PostToFeed(did, text, apiKey); err != nil {
		fmt.Println("Error posting to feed:", err)
		return
	}

	fmt.Println("Posted 'Hello, world' to feed successfully.")
}
