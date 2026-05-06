package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newTestApp() *tview.Application {
	return tview.NewApplication()
}

func TestInputHandlerFieldNotNil(t *testing.T) {
	app := newTestApp()
	h := NewInputHandler(app, nil, nil, nil)
	if h.Field() == nil {
		t.Fatal("expected non-nil input field")
	}
}

func TestInputHandlerQueryEmpty(t *testing.T) {
	app := newTestApp()
	h := NewInputHandler(app, nil, nil, nil)
	if got := h.Query(); got != "" {
		t.Fatalf("expected empty query, got %q", got)
	}
}

func TestInputHandlerClear(t *testing.T) {
	app := newTestApp()
	h := NewInputHandler(app, nil, nil, nil)
	h.inputField.SetText("some text")
	h.Clear()
	if got := h.Query(); got != "" {
		t.Fatalf("expected empty after Clear, got %q", got)
	}
}

func TestInputHandlerOnSearchCallback(t *testing.T) {
	app := newTestApp()
	var captured string
	h := NewInputHandler(app, func(q string) { captured = q }, nil, nil)
	// Simulate text change via SetText which triggers ChangedFunc
	h.inputField.SetText("mysecret")
	if captured != "mysecret" {
		t.Fatalf("expected onSearch to receive 'mysecret', got %q", captured)
	}
}

func TestInputHandlerOnQuitEscape(t *testing.T) {
	app := newTestApp()
	quitCalled := false
	h := NewInputHandler(app, nil, nil, func() { quitCalled = true })

	capture := h.inputField.GetInputCapture()
	if capture == nil {
		t.Skip("no input capture set")
	}
	event := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	capture(event)
	if !quitCalled {
		t.Fatal("expected onQuit to be called on Escape")
	}
}

func TestInputHandlerOnSelectEnter(t *testing.T) {
	app := newTestApp()
	selectCalled := false
	h := NewInputHandler(app, nil, func() { selectCalled = true }, nil)

	capture := h.inputField.GetInputCapture()
	if capture == nil {
		t.Skip("no input capture set")
	}
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	capture(event)
	if !selectCalled {
		t.Fatal("expected onSelect to be called on Enter")
	}
}
