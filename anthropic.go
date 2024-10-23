package computeruse

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

type MessageContent struct {
	Type string `json:"type"` // text, image, tool_use, tool_result

	// for text
	Text string `json:"text,omitempty"`

	// for tool_use
	ID    string        `json:"id,omitempty"`
	Name  string        `json:"name,omitempty"`
	Input *MessageInput `json:"input,omitempty"`

	// for tool_result
	Content   []ContentContent `json:"content,omitempty"`
	ToolUseID string           `json:"tool_use_id,omitempty"`
	IsError   bool             `json:"is_error,omitempty"`
}

type MessageInput struct {
	Action     string `json:"action,omitempty"`
	Text       string `json:"text,omitempty"`
	Coordinate []int  `json:"coordinate,omitempty"`
}

type ContentContent struct {
	Type   string        `json:"type"`
	Source ContentSource `json:"source,omitempty"`
}

type ContentSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type ToollCallResult struct {
	Action     string
	Text       string
	ToolUseID  string
	IsError    bool
	Screenshot []byte // PNG bytes
}

func RunMessages(msgs []Message, width, height int) (Message, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return Message{}, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	data := map[string]interface{}{
		"model":      "claude-3-5-sonnet-20241022",
		"max_tokens": 1024,
		"tools": []map[string]interface{}{
			{
				"type":              "computer_20241022",
				"name":              "computer",
				"display_width_px":  width,
				"display_height_px": height,
				"display_number":    1,
			},
		},
		"messages": msgs,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return Message{}, fmt.Errorf("json marshal: %v", err)
	}

	url := "https://api.anthropic.com/v1/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Message{}, fmt.Errorf("http new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "computer-use-2024-10-22")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Message{}, fmt.Errorf("http client do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return Message{}, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	respMsg := Message{}
	err = json.NewDecoder(resp.Body).Decode(&respMsg)
	if err != nil {
		return Message{}, fmt.Errorf("json decode: %v", err)
	}
	return respMsg, nil
}

func ToolCall(com Computer, msg Message) ([]ToollCallResult, error) {
	results := []ToollCallResult{}
	for _, c := range msg.Content {
		if c.Type != "tool_use" {
			continue
		}
		if c.Name == "computer" {
			switch c.Input.Action {
			case "mouse_move":
				log.Printf("browser: mouse_move %+v", c.Input.Coordinate)
				com.MouseMove(c.Input.Coordinate[0], c.Input.Coordinate[1])
			case "left_click":
				log.Printf("browser: left_click")
				com.LeftClick()
			case "type":
				log.Printf("browser: type %v", c.Input.Text)
				com.Type(c.Input.Text)
			case "key":
				log.Printf("browser: key %v", c.Input.Text)
				com.Key(c.Input.Text)
			case "screenshot":
				log.Printf("browser: screenshot")
			default:
				log.Printf("browser: action %v is not implemented", c.Input.Action)
			}

			screenshot := com.Screenshot()
			results = append(results, ToollCallResult{
				Action:     c.Input.Action,
				Text:       c.Input.Text,
				ToolUseID:  c.ID,
				Screenshot: screenshot,
			})
		}
	}
	return results, nil
}

func NewToolResponseMessage(toolresults []ToollCallResult) Message {
	content := []MessageContent{}
	for _, r := range toolresults {
		content = append(content, MessageContent{
			Type:      "tool_result",
			ToolUseID: r.ToolUseID,
			IsError:   r.IsError,
			Content: []ContentContent{
				{
					Type: "image",
					Source: ContentSource{
						Type:      "base64",
						MediaType: "image/png",
						Data:      base64.StdEncoding.EncodeToString(r.Screenshot),
					},
				},
			},
		})
	}
	return Message{
		Role:    "user",
		Content: content,
	}
}
