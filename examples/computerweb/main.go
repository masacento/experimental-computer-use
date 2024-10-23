package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	computeruse "github.com/masacento/experimental-computer-use"
)

type SSEServer struct {
	clients map[chan string]bool
	mutex   sync.RWMutex
	items   []string
}

func NewSSEServer() *SSEServer {
	return &SSEServer{
		clients: make(map[chan string]bool),
		items:   []string{},
	}
}

func (sse *SSEServer) addClient() chan string {
	sse.mutex.Lock()
	defer sse.mutex.Unlock()

	client := make(chan string, 10)
	sse.clients[client] = true
	return client
}

func (sse *SSEServer) removeClient(client chan string) {
	sse.mutex.Lock()
	defer sse.mutex.Unlock()

	delete(sse.clients, client)
	close(client)
}

func (sse *SSEServer) broadcast(item string) {
	sse.mutex.RLock()
	defer sse.mutex.RUnlock()
	item = strings.TrimSpace(item)
	item = strings.ReplaceAll(item, "\n", "<br/>")

	for client := range sse.clients {
		select {
		case client <- item:
		default:
		}
	}
}

func (sse *SSEServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := sse.addClient()
	defer sse.removeClient(client)

	for _, item := range sse.items {
		fmt.Fprintf(w, "data: %s\n\n", item)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	for {
		select {
		case item := <-client:
			fmt.Fprintf(w, "data: <article>%s</article>\n\n", item)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

func (sse *SSEServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, indexHTML)
}

const indexHTML = `
<!DOCTYPE html>
<html>

<head>
    <title>experimental computer use</title>
    <script src="https://unpkg.com/htmx.org@2.0.3"></script>
    <script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
</head>

<body>
    <script>
        document.body.addEventListener('htmx:afterSwap', (event) => {
            setTimeout(() => {
                const mdiv = document.getElementById('messages')
                window.scrollTo({
                    top: mdiv.scrollHeight,
                    behavior: 'smooth'
                })
            }, 100)
        });
    </script>
    <main class="container">
        <div class="list-container">
            <h2>experimental browser use</h2>
            <form method="post" hx-post="/run" hx-swap="none">
				<div style="display: flex; align-items: center;">
                	<input type="text" name="prompt" value="find cute cats">
                	<button style="margin-left: 1rem; width: 10rem;" type="submit">Run</button>
				</div>
            </form>

            <div id="messages" hx-ext="sse" sse-connect="/events" sse-swap="message" hx-swap="beforeend" />

        </div>
    </main>
</body>

</html>
`

func (sse *SSEServer) handleRun(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("prompt")
	if prompt == "" {
		http.Error(w, "Prompt cannot be empty", http.StatusBadRequest)
		return
	}

	html := fmt.Sprintf("User: %s", prompt)
	sse.items = append(sse.items, html)
	sse.broadcast(html)

	top := 100
	left := 5
	width := 300
	height := 600
	maxturn := 20
	rb := computeruse.NewRobot(computeruse.Rect{top, left, width, height})

	msg := computeruse.Message{
		Role: "user",
		Content: []computeruse.MessageContent{
			{
				Type: "text",
				Text: prompt,
			},
		},
	}
	sse.broadcast(prompt)

	msgs := []computeruse.Message{msg}

	for i := 0; i < maxturn; i++ {
		resp, err := computeruse.RunMessages(msgs, width, height)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v\n", resp)
		respmes := resp.Content[0].Text
		if len(resp.Content) == 1 && resp.Content[0].Type == "text" {
			html := fmt.Sprintf("Assistant: %s", respmes)
			sse.items = append(sse.items, html)
			sse.broadcast(html)
			break
		}
		msgs = append(msgs, resp)

		result, err := computeruse.ToolCall(rb, resp)
		if err != nil {
			log.Fatal(err)
		}
		if len(result) == 0 {
			continue
		}

		resptool := result[len(result)-1]
		html := fmt.Sprintf("Assistant: %s [ToolCall: %s %s]<br/><img src=\"data:image/png;base64,%s\"/>", respmes, resptool.Action, resptool.Text, base64.StdEncoding.EncodeToString(resptool.Screenshot))
		sse.items = append(sse.items, html)
		sse.items = append(sse.items, "Done")
		sse.broadcast(html)
		msgs = append(msgs, computeruse.NewToolResponseMessage(result))
	}
}

func main() {
	sse := NewSSEServer()

	http.HandleFunc("/", sse.handleIndex)
	http.HandleFunc("/run", sse.handleRun)
	http.HandleFunc("/events", sse.handleSSE)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
