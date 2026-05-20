package main

import (
	"fmt"
	"os"
	"context"
	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"
)

func main() {
	fmt.Println("test")
	
	getConfigPath()
	
	// anytypeAuth()

}

func getConfigPath() {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	} 
	fmt.Println(dir)
	return dir
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
