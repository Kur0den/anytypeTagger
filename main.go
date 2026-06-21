package main

import (
	"context"
	"fmt"
	"os"

	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"

	"github.com/openai/openai-go/v3"
)


type Config struct {
	LLM struct {
		Host			string 	`json:"host"`
		Port			string 	`json:"port"`
		IsHttps		bool		`json:"isHttps"`
		Endpoint	string 	`json:"endpoint"`
		Key				string 	`json:"key"`
	} `json:"llm"`
}


func main() {
	// configを取得
	config, err := getConfigPath()
	
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)
	// anytypeAuth()
}

func getConfigPath() (dir string, err error) {
	// configPathを取得
	dir, err = os.UserConfigDir()
	if err != nil {
		return 
	}
	fmt.Println(dir)
	// 保存するディレクトリを定義
	configPath := dir + "/AnytypeTagger"
	// ディレクトリが存在するか確認
	_, err = os.Stat(configPath)
	if err != nil {
		// 存在しない場合は新規作成
		fmt.Println("configディレクトリが存在しないため新規作成します")
		err = os.Mkdir(configPath, 0755)
		if err != nil {
			// 作成失敗時はreturn
			return
		}
		_, err = os.Stat(configPath)
		if err != nil {
			return
		}
	}
	return configPath, err
}

func getConfig() (Config, error) {
	config := Config {}
	
	config.LLM.Host = "192.168.0.20"
	config.LLM.Port = "8080"
	config.LLM.IsHttps = false
	config.LLM.Endpoint = "http://192.168.0.20:8080/v1"

	return config, nil
}


func anytypeAuth() {
	// clientを定義
	client := anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
	)
	
	// Anytypeの認証を呼び出し
	ctx := context.Background()
	auth, _ := client.Auth().CreateChallenge(ctx, "AnytypeTagger")
	
	fmt.Print("code: ")
	var code string
	fmt.Scanln(&code)

	// 入力されたコードで問い合わせ
	token, err := client.Auth().CreateApiKey(ctx, auth.ChallengeID, code)
	
	fmt.Println(token)
	fmt.Println(err)
	
	// 取得したAPIキーを元にclientを再定義
	client = anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
		anytype.WithAppKey(token.ApiKey),
	)

	fmt.Println(client)
}
