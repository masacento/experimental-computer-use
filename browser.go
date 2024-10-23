package computeruse

import (
	"log"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

type Browser struct {
	browser *rod.Browser
	page    *rod.Page
}

func NewBrowser() *Browser {
	browser := rod.New().MustConnect()
	return &Browser{browser, nil}
}

func (b *Browser) Open(url string) {
	log.Println("open", url)
	page, err := b.browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		log.Println("Error opening page:", err)
		return
	}
	page.MustWaitStable()
	b.page = page
}

func (b *Browser) Screenshot() []byte {
	screenshot, err := b.page.Screenshot(false, nil)
	if err != nil {
		log.Println("Error taking screenshot:", err)
		return nil
	}
	return screenshot
}

func (b *Browser) CursorPosition() (x, y int) {
	mouse := b.page.Mouse
	return int(mouse.Position().X), int(mouse.Position().Y)
}

func (b *Browser) Key(key string) {
	keyb := b.page.Keyboard
	switch strings.ToLower(key) {
	case "enter":
		keyb.Press(input.Enter)
	case "return":
		keyb.Press(input.Enter)
	case "delete":
		keyb.Press(input.Delete)
	case "tab":
		keyb.Press(input.Tab)
	case "escape":
		keyb.Press(input.Escape)
	case "left":
		keyb.Press(input.ArrowLeft)
	case "right":
		keyb.Press(input.ArrowRight)
	case "up":
		keyb.Press(input.ArrowUp)
	case "down":
		keyb.Press(input.ArrowDown)
	case "page_up":
		keyb.Press(input.PageUp)
	case "page_down":
		keyb.Press(input.PageDown)
	default:
		log.Printf("key: %v is not implemented", key)
	}
	b.page.MustWaitStable()
}

func (b *Browser) Type(text string) {
	page := b.page
	page.InsertText(text)
}

func (b *Browser) MouseMove(x, y int) {
	mouse := b.page.Mouse
	mouse.MustMoveTo(float64(x), float64(y))
}

func (b *Browser) LeftClick() {
	mouse := b.page.Mouse
	mouse.MustDown("left")
	mouse.MustUp("left")
	b.page.MustWaitStable()
}

func (b *Browser) RightClick() {
	mouse := b.page.Mouse
	mouse.MustDown("right")
	mouse.MustUp("right")
	b.page.MustWaitStable()
}
