package computeruse

// https://docs.anthropic.com/en/docs/build-with-claude/computer-use#computer-tool
// * `key`: Press a key or key-combination on the keyboard.
// 	- This supports xdotool's `key` syntax.
// 	- Examples: "a", "Return", "alt+Tab", "ctrl+s", "Up", "KP_0" (for the numpad 0 key).
// * `type`: Type a string of text on the keyboard.
// * `cursor_position`: Get the current (x, y) pixel coordinate of the cursor on the screen.
// * `mouse_move`: Move the cursor to a specified (x, y) pixel coordinate on the screen.
// * `left_click`: Click the left mouse button.
// * `left_click_drag`: Click and drag the cursor to a specified (x, y) pixel coordinate on the screen.
// * `right_click`: Click the right mouse button.
// * `middle_click`: Click the middle mouse button.
// * `double_click`: Double-click the left mouse button.
// * `screenshot`: Take a screenshot of the screen.""",

type Computer interface {
	MouseMove(x, y int)
	LeftClick()
	RightClick()
	Type(text string)
	Key(key string)
	Screenshot() []byte
	CursorPosition() (x int, y int)
}
