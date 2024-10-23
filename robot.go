package computeruse

import (
	"bytes"
	"image/png"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/nfnt/resize"
)

func NewRobot(bounds Rect) Computer {
	return &Robot{bounds}
}

type Robot struct {
	Bounds Rect
}

type Rect struct {
	Top    int
	Left   int
	Width  int
	Height int
}

const sleepMilli = 1000

func (r *Robot) MouseMove(x, y int) {
	x += r.Bounds.Left
	y += r.Bounds.Top
	if x > r.Bounds.Width {
		x = r.Bounds.Width
	}
	if y > r.Bounds.Height {
		y = r.Bounds.Height
	}

	robotgo.Move(x, y)
}

func (r *Robot) LeftClick() {
	robotgo.Click("left")
	robotgo.MilliSleep(sleepMilli)
}

func (r *Robot) RightClick() {
	robotgo.Click("right")
	robotgo.MilliSleep(sleepMilli)
}

func (r *Robot) Type(text string) {
	robotgo.TypeStr(text)
	robotgo.MilliSleep(sleepMilli)
}

func (r *Robot) Key(key string) {
	key = strings.ToLower(key)
	robotgo.KeyTap(key)
	robotgo.MilliSleep(sleepMilli)
}

func (r *Robot) Screenshot() []byte {
	bit := robotgo.CaptureScreen(r.Bounds.Left, r.Bounds.Top, r.Bounds.Width, r.Bounds.Height)
	defer robotgo.FreeBitmap(bit)

	img := robotgo.ToImage(bit)
	resized := resize.Resize(uint(r.Bounds.Width), uint(r.Bounds.Height), img, resize.Lanczos3)

	w := new(bytes.Buffer)
	err := png.Encode(w, resized)
	if err != nil {
		return nil
	}
	return w.Bytes()
}

func (r *Robot) CursorPosition() (x int, y int) {
	x, y = robotgo.Location()
	return x, y
}
