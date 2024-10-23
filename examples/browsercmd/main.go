package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	computeruse "github.com/masacento/experimental-computer-use"
)

func main() {
	br := computeruse.NewBrowser()
	br.Open("https://www.duckduckgo.com/")

	prompt := "search for images of calico cats"

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
		resp, err := computeruse.RunMessages(msgs, 1024, 768)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.Content) == 1 && resp.Content[0].Type == "text" {
			fmt.Println("finish:", resp.Content[0].Text)
			break
		}
		msgs = append(msgs, resp)
		fmt.Printf("%+v\n", resp)
		result, err := computeruse.ToolCall(br, resp)
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
