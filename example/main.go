package main

import "fmt"

func main() {
	handle := "totototo"
	appPassword := "myapppassword"
	actor := "tototot"
	limit := 4
	text := "Hello, world"

	did, err := ResolveDID(handle)
	if err != nil {
		fmt.Println("Error resolving DID:", err)
		return
	}

	apiKey, err := GetApiKey(did, appPassword)
	if err != nil {
		fmt.Println("Error getting API key:", err)
		return
	}

	userFeed, err := GetUserFeed(actor, limit, apiKey)
	if err != nil {
		fmt.Println("Error getting user feed:", err)
		return
	}

	fmt.Println("User feed:", userFeed)

	if _, err := PostToFeed(did, text, apiKey); err != nil {
		fmt.Println("Error posting to feed:", err)
		return
	}

	fmt.Println("Posted 'Hello, world' to feed successfully.")
}
