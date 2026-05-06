package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InputHandler manages the search input field and keyboard event routing.
type InputHandler struct {
	app       *tview.Application
	inputField *tview.InputField
	onSearch  func(query string)
	onQuit    func()
	onSelect  func()
}

// NewInputHandler creates an InputHandler wired to the provided callbacks.
func NewInputHandler(app *tview.Application, onSearch func(string), onSelect func(), onQuit func()) *InputHandler {
	h := &InputHandler{
		app:      app,
		onSearch: onSearch,
		onSelect: onSelect,
		onQuit:   onQuit,
	}

	field := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDefault)

	field.SetChangedFunc(func(text string) {
		if h.onSearch != nil {
			h.onSearch(strings.TrimSpace(text))
		}
	})

	field.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			if h.onQuit != nil {
				h.onQuit()
			}
			return nil
		case tcell.KeyEnter:
			if h.onSelect != nil {
				h.onSelect()
			}
			return nil
		}
		return event
	})

	h.inputField = field
	return h
}

// Field returns the underlying tview InputField for layout composition.
func (h *InputHandler) Field() *tview.InputField {
	return h.inputField
}

// Query returns the current trimmed text in the input field.
func (h *InputHandler) Query() string {
	return strings.TrimSpace(h.inputField.GetText())
}

// Clear resets the input field text.
func (h *InputHandler) Clear() {
	h.inputField.SetText("")
}
