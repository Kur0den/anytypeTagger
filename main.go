package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"errors"

	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	// configを取得
	configPath, err := getConfigPath()
	
	if err != nil {
		log.Fatalf("configPathがしゅとくできませんでした: %v", err)
	}

	config, err := getConfig(configPath)
	
	if err != nil {
		log.Fatalf("configがしゅとくできませんでした: %v", err)
	}

	app := &Application{
		config: config,
	}

	app.Run()

}

func getConfigPath() (string, error) {
	// configPathを取得
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("os.UserConfigDirの取得に失敗したよ")
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
			return "", fmt.Errorf("ディレクトリの作成に失敗したよ")
		}
		_, err = os.Stat(configPath)
		if err != nil {
			return "", fmt.Errorf("ディレクトリの確認に失敗したよ")
		}
	}
	return configPath, err
}

func getConfig(configPath string) (Config, error) {

	// _ -> file
	_, err := os.Open(configPath + "/config.json")
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return Config{}, fmt.Errorf("権限不足で設定ファイルをひらけなかったよ")
		} else if errors.Is(err, os.ErrNotExist) {
			fmt.Println("設定ファイルが見つからなかったよ あたらしくつくるね")
		}
	}
	fmt.Println("")
	var config Config

	if true {
		config.LLM.Endpoint = "http://192.168.0.20:8081/v1"
		config.LLM.Token = "DUMMY"
		config.Anytype.Endpoint = "http://localhost:31009"
		config.Anytype.TagId = ""
	} else {
		//TODO fileを読む
	}
	return config, nil

}

func (app *Application) Run() {
	
	ctx := context.Background()

	err := app.anytypeAuth(ctx)
	if err != nil {
		log.Fatalf("認証がうまくいかなかったよ: %v", err)
	}

	err = app.getAnytypeObjId()
	if err != nil {
		log.Fatalf("urlが取得できなかったよ: %v", err)
	}

	content, err := app.getAnytypeObjMd(ctx)
	if err != nil {
		log.Fatalf("対象のオブジェクトが取得できなかったよ: %v", err)
	}
	fmt.Println(content)

	tags, err := app.getAnytypeTags(ctx)
	if err != nil {
		log.Fatalf("オブジェクトが存在するスペースからのタグ取得に失敗したよ: %v", err)
	}



	llmClient := openai.NewClient(
		option.WithAPIKey(app.config.LLM.Token),
		option.WithBaseURL(app.config.LLM.Endpoint),
		)

	models, err := llmClient.Models.List(context.Background())

	if err != nil {
		fmt.Println("modelsが取得できなかったよ")
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





func (app *Application) anytypeAuth(ctx context.Context) (error){
	
	var token string
	
	// configにtokenがない場合はチャレンジを実行
	if app.config.Anytype.Token != "" {
		// clientを定義
		client := anytype.NewClient(
			anytype.WithBaseURL(app.config.Anytype.Token),
		)
		
		// Anytypeの認証を呼び出し
		auth, _ := client.Auth().CreateChallenge(ctx, "AnytypeTagger")
		
		fmt.Print("code: ")
		var code string
		fmt.Scanln(&code)

		// 入力されたコードで問い合わせ
		receivedToken, err := client.Auth().CreateApiKey(ctx, auth.ChallengeID, code)
		
		if err != nil {
			return fmt.Errorf("認証のチャレンジにしっぱいしたかも: %w", err)
		}

		token = receivedToken.ApiKey
	} else { // configにすでにtokenがある場合はそれを使用
		token = app.config.Anytype.Token
	}
	
	// 取得したAPIキーを元にclientを再定義
	client := anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
		anytype.WithAppKey(token),
	)

	app.Anytype.client = client

	return  nil
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


func (app *Application) getAnytypeObjId() (error) {
	inputUrl, err := readStdio("タグ付けをしたいAnytypeのObjのディープリンクをいれてね")
	if err != nil {
		return err
	}
	url, err := url.Parse(inputUrl)
	if err != nil {
		return fmt.Errorf("urlのパースができなかったよ, %v", err)
	}
	if url.Scheme != "anytype" {
		return fmt.Errorf("anytypeのスキーマじゃないみたいだよ")
	} else if url.Host != "object" {
		return fmt.Errorf("objectに対するディープリンクじゃないみたいだよ")
	}

	objIds := AnytypeObjIds {}

	objIds.ObjectId = url.Query().Get("objectId")
	objIds.SpaceId = url.Query().Get("spaceId")
	
	if objIds.ObjectId == "" || objIds.SpaceId == "" {
		return fmt.Errorf("objectIdかspaceIdがディープリンクに存在しないみたいだよ")
	}

	app.Anytype.objIds = objIds

	return nil
}

func (app *Application) getAnytypeObjMd(ctx context.Context) (string ,error){
	targetobj, err := app.Anytype.client.Space(app.Anytype.objIds.SpaceId).Object(app.Anytype.objIds.ObjectId).Get(ctx)
	if err != nil {
		return "", fmt.Errorf("objectの取得にしっぱいしたよ: %v", err)
	}
	return targetobj.Object.Markdown, nil
}

func (app *Application) getAnytypeTags(ctx context.Context) ([]anytype.Tag, error){
	tags, err := app.Anytype.client.Space(app.Anytype.objIds.SpaceId).Property(app.config.Anytype.TagId).Tags().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("Tagsの取得に失敗したよ: %v", err)
	}

	return tags, nil
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
		return nil, fmt.Errorf("LLMへの問い合わせに失敗したよ: %v", err)
	}
	return resp, nil
}
