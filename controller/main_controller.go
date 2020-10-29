package controller

import (
	"strings"

	"github.com/thomgray/notebee/constants"
	"github.com/thomgray/notebee/model"

	"github.com/thomgray/egg"
	"github.com/thomgray/notebee/config"
	"github.com/thomgray/notebee/view"
)

type inputCommand struct {
	key        egg.Key
	activeMode constants.ActiveMode
}

// MainController ...
type MainController struct {
	View              *view.MainView
	InputView         *view.InputView
	ModalMenu         *view.ModalMenu
	CompletionView    *view.CompletionView
	SearchResultsView *view.SearchResultsView
	Config            *config.Config
	FileManager       *model.FileManager
	activeDocument    *model.Document
	activeFile        *model.File
	lastCommand       inputCommand
	activeMode        constants.ActiveMode
}

// Mode ...
type Mode uint8

// Mode ...
const (
	ModeInput Mode = iota
	ModeMenu
)

var app *egg.Application
var mode Mode = ModeInput
var inputMode constants.InputMode = constants.InputModeTraverse

// InitMainController ...
func InitMainController(config *config.Config) *MainController {
	app = egg.InitOrPanic()

	mc := MainController{
		View:              view.MakeMainView(app),
		InputView:         view.MakeInputView(app),
		ModalMenu:         view.MakeModalMenu(),
		CompletionView:    view.MakeCompletionView(),
		SearchResultsView: view.SearchResultsView{}.New(),
		Config:            config,
		FileManager:       model.MakeFileManager(config),
	}

	app.OnResizeEvent(func(re *egg.ResizeEvent) {
		mc.View.Refit(re.Width, re.Height)
		mc.SearchResultsView.Refit(re.Width, re.Height)
		app.ReDraw()
	})

	// app.AddViewController(mc.ModalMenu)
	app.AddView(mc.SearchResultsView.View)

	app.OnKeyEvent(mc.keyEventDelegate)

	mc.init()

	return &mc
}

func (mc *MainController) init() {
	mc.reloadFiles()
	bootstrapCommands()
}

func (mc *MainController) reloadFiles() {
	// mc.FileManager.LoadFiles(mc.Config.NotePaths)
}

func (mc *MainController) setMode(mode constants.ActiveMode) {
	mc.activeMode = mode
	switch mode {
	case constants.ActiveModeDefault:
		mc.CompletionView.Close()
		mc.SearchResultsView.Close()
	case constants.ActiveModeAutocomplete:
		mc.SearchResultsView.Close()
	case constants.ActiveModeSearchResultSelect:
		mc.CompletionView.Close()
	}
}

func (mc *MainController) keyEventDelegate(e *egg.KeyEvent) {
	defer app.ReDraw()
	switch e.Key {
	case egg.KeyEsc:
		mc.setMode(constants.ActiveModeDefault)
		e.SetPropagate(false)
		app.ReDraw()
		return
	}

	defer (func() {
		mc.lastCommand.activeMode = mc.activeMode
		mc.lastCommand.key = e.Key
	})()

	switch mc.activeMode {
	case constants.ActiveModeDefault:
		// if you are entering tab mode
		if e.Key == egg.KeyTab && mc.lastCommand.key == egg.KeyTab && mc.CompletionView.IsOpen() {
			// log.Println("In tabby mode")
			mc.setMode(constants.ActiveModeAutocomplete)
			mc.handleCompltionModeEvent(e)
		} else {
			mc.handleEventInputMode(e)
		}
	case constants.ActiveModeAutocomplete:
		mc.handleCompltionModeEvent(e)
	case constants.ActiveModeSearchResultSelect:
		mc.handleSearchResultModeEvent(e)
	}
}

func (mc *MainController) handleEventInputMode(e *egg.KeyEvent) {
	if mc.InputView.GetCursorX() == 0 {
		switch e.Char {
		case '?':
			e.SetPropagate(false)
			mc.setInputMode(constants.InputModeSearch)
			app.ReDraw()
			return
		case '>':
			e.SetPropagate(false)
			mc.setInputMode(constants.InputModeTraverse)
			app.ReDraw()
			return
		case ':':
			e.SetPropagate(false)
			mc.setInputMode(constants.InputModeCommand)
			app.ReDraw()
			return
		}
	}

	switch e.Key {
	case egg.KeyEnter:
		mc.handleEnter(e)
	case egg.KeyTab:
		e.SetPropagate(false)
		mc.handleAutocomplete(mc.InputView.GetTextContentString())
	case egg.KeyUp, egg.KeyDown:
		e.SetPropagate(false)
		mc.CompletionView.SetVisible(false)
		mc.View.HandleKeyEvent(e)
	}
}

func (mc *MainController) handleEventMenuMode(e *egg.KeyEvent) {
	switch e.Char {
	case 'x':
		app.Stop()
	}
}

func (mc *MainController) toggleMode() Mode {
	if mode == ModeInput {
		app.SetFocusedView(nil)
		mode = ModeMenu
	} else {
		mc.InputView.GainFocus()
		mode = ModeInput
	}
	return mode
}

func (mc *MainController) setInputMode(m constants.InputMode) {
	if inputMode == m {
		return
	}
	// old := mode
	inputMode = m
	mc.InputView.SetMode(inputMode)
}

func (mc *MainController) handleEnter(e *egg.KeyEvent) {
	e.SetPropagate(false)
	mc.CompletionView.Close()
	txt := mc.InputView.GetTextContentString()
	switch inputMode {
	case constants.InputModeTraverse:
		// if !mc.handleSpecial(txt) {
		mc.handleTraverse(txt)
		// }
	case constants.InputModeSearch:
		mc.handleSearch(txt)
	case constants.InputModeCommand:
		mc.handleCommand(txt)
	}
}

// SetActiveFile ...
func (mc *MainController) SetActiveFile(f *model.File) {
	mc.activeFile = f
	mc.View.SetActiveFile(f)
	mc.InputView.SetTextContentString("")
	mc.InputView.SetCursorX(0)
}

func matchesLocation(f *model.File, path string) bool {
	for _, loc := range f.Locations {
		if strings.EqualFold(loc.RelativePathWithName, path) {
			// path matches fully
			return true
		}
	}
	return false
}

func (mc *MainController) handleAutocomplete(str string) {
	switch inputMode {
	case constants.InputModeTraverse:
		mc.handleAutocompleteNote(str)
	}
}

// Start ...
func (mc *MainController) Start() {
	defer app.Start()
}
