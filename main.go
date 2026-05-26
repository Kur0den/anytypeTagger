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
	info, err := os.Stat(configPath)
	if err != nil {
		// 存在しない場合は新規作成
		fmt.Println("configディレクトリが存在しないため新規作成します")
		err = os.Mkdir(configPath, 0755)
		if err != nil {
			// 作成失敗時はreturn
			return
		}
	}
	fmt.Println(info)
	return 
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
