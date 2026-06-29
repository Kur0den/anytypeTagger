package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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
	_, err := getConfigPath()
	
	if err != nil {
		log.Fatalf("configPathがしゅとくできませんでした: %v", err)
	}

	config, err := getConfig()
	
	if err != nil {
		log.Fatalf("configがしゅとくできませんでした: %v", err)
	}

	// client, err := anytypeAuth()
	// if err != nil {
	//   log.Fatalf("にんしょうがうまくいかなかったよ: %v", err)
	// }

	url, err := getAnytypeObjUrl()
	if err != nil {
		log.Fatalf("urlがしゅとくできませんでした: %v", err)
	}
	fmt.Println(url)


	client := openai.NewClient(
		option.WithAPIKey(config.LLM.Key),
		option.WithBaseURL(config.LLM.Endpoint),
	)

	models, err := client.Models.List(context.Background())

	if err != nil {
		fmt.Println("modelsがしゅとくできませんでした")
		panic(err)
	}

	modelId := ""
	for _, model := range models.Data {
		fmt.Println("Model:", model.ID)
		modelId = model.ID
	}

	resp, err := getLLMResponse(modelId)


	fmt.Println(resp.Choices[0].Message.Content)
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
	config.LLM.Endpoint = "http://192.168.0.20:8081/v1"
	config.LLM.Key = "DUMMY"
	return config, nil

}


func anytypeAuth() (*anytype.Client, error){
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
	
	if err != nil {
		return nil, fmt.Errorf("にんしょうのちゃれんじにしっぱいしたかも: %w", err)
	}

	// 取得したAPIキーを元にclientを再定義
	client = anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
		anytype.WithAppKey(token.ApiKey),
	)

	return client, nil
}

func readStdio(prompt) (string, err){
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("stdinのよみとりができなかったよ")
	}
	return text, nil
}


func getAnytypeObjUrl() (string, error) {
	inputUrl, err = readStdio("タグ付けをしたいAnytypeのObjのディープリンクをいれてね")
	
}

func getLLMResponse(modelId) {
	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: modelId,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage("あなたは可愛い幼女なAIアシスタントです"),
			openai.UserMessage("こんにちは"),
		},
	})

	if err != nil {
		fmt.Println("問い合わせに失敗したかも")
		panic(err)
	}
}
