package main

import (
	"fmt"
	"context"
	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"
)

func main() {
	fmt.Println("test")
	
	
	anytypeAuth()
}


func anytypeAuth() {
	client := anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
	)
	
	ctx := context.Background()
	auth, _ := client.Auth().CreateChallenge(ctx, "AnytypeTagger")
	
	fmt.Print("code: ")
	var code string
	fmt.Scanln(&code)

	token, err := client.Auth().CreateApiKey(ctx, auth.ChallengeID, code)
	
	fmt.Println(token)
	fmt.Println(err)

	client = anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
		anytype.WithAppKey(token.ApiKey),
	)

	fmt.Println(client)
}
