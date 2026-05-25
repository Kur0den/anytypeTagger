package main

import (
	"context"
	"fmt"
	"os"

	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"
)

func main() {
	fmt.Println("test")
	
	config, err := getConfigPath()
	
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)
	// anytypeAuth()

}

func getConfigPath() (dir string, err error) {
	dir, err = os.UserConfigDir()
	if err != nil {
		return 
	}
	fmt.Println(dir)
	configPath := dir + "/AnytypeTagger"
	info, err := os.Stat(configPath)
	if err != nil {
		fmt.Println("configディレクトリが存在しないため新規作成します")
		err = os.Mkdir(configPath, 0755)
		if err != nil {
			return
		}
	}
	fmt.Println(info)
	return 
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
