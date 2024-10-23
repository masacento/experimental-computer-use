package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	computeruse "github.com/masacento/experimental-computer-use"
)

func main() {
	width := 1024
	height := 40 // limit to menu bar
	rb := computeruse.NewRobot(computeruse.Rect{0, 0, width, height})

	prompt := "what app is running on this mac? click that app menu."

	msg := computeruse.Message{
		Role: "user",
		Content: []computeruse.MessageContent{
			{
				Type: "text",
				Text: prompt,
			},
		},
	}

	msgs := []computeruse.Message{msg}

	for i := 0; i < 10; i++ {
		resp, err := computeruse.RunMessages(msgs, width, height)
		if err != nil {
			log.Fatal(err)
		}
		msgs = append(msgs, resp)
		fmt.Printf("%+v\n", resp)
		if len(resp.Content) == 1 && resp.Content[0].Type == "text" {
			fmt.Println("finish:", resp.Content[0].Text)
			break
		}
		result, err := computeruse.ToolCall(rb, resp)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range result {
			os.MkdirAll("_screenshots", 0755)
			os.WriteFile(filepath.Join("_screenshots", r.ToolUseID+".png"), r.Screenshot, 0644)
		}
		msgs = append(msgs, computeruse.NewToolResponseMessage(result))
	}
}
