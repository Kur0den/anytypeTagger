package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

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

	anytypeClient, err := anytypeAuth()
	if err != nil {
	  log.Fatalf("にんしょうがうまくいかなかったよ: %v", err)
	}

	objIds, err := getAnytypeObjId()
	if err != nil {
		log.Fatalf("urlがしゅとくできませんでした: %v", err)
	}
	
	content, err := getAnytypeObjMd(context.Background(), *anytypeClient, objIds)
	if err != nil {
		log.Fatalf("たいしょうのおぶじぇくとがしゅとくできませんでした: %v", err)
	}
	fmt.Println(content)
	return

	llmClient := openai.NewClient(
		option.WithAPIKey(config.LLM.Token),
		option.WithBaseURL(config.LLM.Endpoint),
	)

	models, err := llmClient.Models.List(context.Background())

	if err != nil {
		fmt.Println("modelsがしゅとくできませんでした")
		panic(err)
	}

	modelId := ""
	for _, model := range models.Data {
		fmt.Println("Model:", model.ID)
		modelId = model.ID
	}

	resp, err := getLLMResponse(modelId, llmClient)

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
	config.LLM.Endpoint = "http://192.168.0.20:8081/v1"
	config.LLM.Token = "DUMMY"
	config.Anytype.Endpoint = "http://localhost:31009"
	config.Anytype.TagId = ""
	return config, nil
}


func anytypeAuth(config Config) (*anytype.Client, error){
	
	var token string

	if config.Anytype.Token != "" {
		// clientを定義
		client := anytype.NewClient(
			anytype.WithBaseURL(config.Anytype.Token),
		)
		
		// Anytypeの認証を呼び出し
		ctx := context.Background()
		auth, _ := client.Auth().CreateChallenge(ctx, "AnytypeTagger")
		
		fmt.Print("code: ")
		var code string
		fmt.Scanln(&code)

		// 入力されたコードで問い合わせ
		receivedToken, err := client.Auth().CreateApiKey(ctx, auth.ChallengeID, code)
		
		if err != nil {
			return nil, fmt.Errorf("にんしょうのチャレンジしっぱいしたかも: %w", err)
		}

		// 取得したAPIキーを元にclientを再定義
		token = receivedToken.ApiKey
	} else {
		token = config.Anytype.Token
	}

	client := anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
		anytype.WithAppKey(token),
	)

	return &client, nil
}

func readStdio(prompt string) (string, error){
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + ": ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("stdinのよみとりができなかったよ")
	}
	text = strings.TrimSpace(text)
	return text, nil
}


func getAnytypeObjId() (*AnytypeObjIds, error) {
	inputUrl, err := readStdio("タグ付けをしたいAnytypeのObjのディープリンクをいれてね")
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(inputUrl)
	if err != nil {
		return nil, fmt.Errorf("urlのパースができなかったよ, %v", err)
	}
	if url.Scheme != "anytype" {
		return nil, fmt.Errorf("anytypeのスキーマじゃないみたいだよ")
	} else if url.Host != "object" {
		return nil, fmt.Errorf("objectに対するディープリンクじゃないみたいだよ")
	}

	objIds := AnytypeObjIds {}

	objIds.ObjectId = url.Query().Get("objectId")
	objIds.SpaceId = url.Query().Get("spaceId")
	
	if objIds.ObjectId == "" || objIds.SpaceId == "" {
		return nil, fmt.Errorf("objectIdかspaceIdがディープリンクに存在しないみたいだよ")
	}

	return &objIds, nil
}

func getAnytypeObjMd(ctx context.Context, client anytype.Client, objIds *AnytypeObjIds) (string ,error){
	targetobj, err := client.Space(objIds.SpaceId).Object(objIds.ObjectId).Get(ctx)
	if err != nil {
		return "", fmt.Errorf("objectの取得にしっぱいしたよ: %v", err)
	}
	return targetobj.Object.Markdown, nil
}

func getLLMResponse(modelId string, llmClient openai.Client) (*openai.ChatCompletion, error){
	resp, err := llmClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
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
	return resp, nil
}
