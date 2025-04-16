package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
	"github.com/gin-gonic/gin"
)

type DataForAI struct {
	UserName        string `json:"my_name"`
	UserInfo        string `json:"my_info"`
	CharacterName   string `json:"character_name"`
	CharacterInfo   string `json:"character_info"`
	CharacterGender string `json:"character_gender"`
	WorldView       string `json:"world_view"`
}

type ResponseOfAI struct {
	Err     error  `json:"error"`
	Message string `json:"ai_message"`
}

// anthropic ㅆ련들 ㅈㄴ 많이 바꿔놨네.
// 알파 버전일 때 만들어 놓은거 갖고 와서 코드 수정하는데 1시간 걸림 ㅅㅂ
// 아니 왜 golang만 지원 잘 안해주냐. 아직까지 베타인게 말이 돼??
func GetResponseOfAI(c *gin.Context, data *DataForAI) ResponseOfAI {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "streaming unsupported")
		return ResponseOfAI{Err: errors.New("streaming unsupported"), Message: ""}
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json 변환 실패", err)
		return ResponseOfAI{Err: errors.New("json 변환 실패"), Message: ""}
	}
	systemText := []anthropic.TextBlockParam{
		{
			Type: constant.Text("text"),
			Text: `비언어적 표현과 말을 평균 300자 이상으로 최대한 자세히 작성해주세요.  
비언어적 표현은 앞뒤에 * (별표 기호)를 붙여 감싸며, 말과 자연스럽게 동시에 일어나는 행동은 대사와 함께 써주세요.  
대사 없이 단독 행동은 *...다* 형태로 끝내주세요.  
감정, 표정, 눈빛, 몸짓을 시각적·감각적으로 묘사해 주세요.  
모든 비언어적 행동 표현은 반드시 "~하고 있다", "~이다", "~한다" 등 서술형 문체를 사용해주세요.  
"~하고 있습니다", "~입니다" 같은 존댓말 문체는 행동 표현에서는 절대 사용하지 마세요. 
대사는 자유롭게 작성해도 됩니다.`,
		},
		{
			Type: constant.Text("text"),
			Text: fmt.Sprintf("이 대화에서 '%s'는 사용자(User)이며, 너는 '%s'라는 캐릭터다. 너는 이제부터 %s로서 대화해야 하며, 절대 이 역할을 벗어나지 마라.", data.UserName, data.CharacterName, data.CharacterName) + string(jsonBytes),
		},
		{
			Type: constant.Text("text"),
			Text: string(jsonBytes),
		},
	}

	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 1024,
		System:    systemText,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("")),
		},
	})
	defer stream.Close()

	var finalText strings.Builder

	for stream.Next() {
		event := stream.Current()

		if event.Type == "content_block_delta" {
			// event.Delta는 구조체로 정의되어 있으므로, 직접 필드에 접근합니다.
			if event.Delta.Type == "text_delta" {
				text := event.Delta.Text
				finalText.WriteString(text)
				fmt.Fprintf(c.Writer, "%s", text)
				flusher.Flush()
			}
		}
	}

	if err := stream.Err(); err != nil {
		c.String(http.StatusInternalServerError, "Stream error: "+err.Error())
		return ResponseOfAI{Err: err, Message: ""}
	}

	return ResponseOfAI{Err: nil, Message: ""}
}
