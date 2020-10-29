package controller

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/thomgray/egg"
	"github.com/thomgray/notebee/constants"
	"github.com/thomgray/notebee/model"
)

func (mc *MainController) handleAutocompleteNote(str string) {
	completeSuggestions := mc.suggestAutocompletions(str)

	if len(completeSuggestions) == 1 {
		newQuery := completeSuggestions[0].CompletionStr()
		mc.InputView.SetTextContentString(newQuery)
		mc.InputView.SetCursorX(runewidth.StringWidth(newQuery))
	} else {
		mc.CompletionView.SetCompletions(completeSuggestions)
		mc.CompletionView.Open()
	}
}

func (mc *MainController) suggestAutocompletions(fragment string) []model.AutocompleteResult {
	res := make([]model.AutocompleteResult, 0)
	allFiles := mc.FileManager.FindSupportedFilePaths()
	topCompleteDir := filepath.Dir(fragment)

	var resContains func(string) bool
	resContains = func(str string) bool {
		for _, s := range res {
			if s.Str == str {
				return true
			}
		}
		return false
	}

	for _, f := range allFiles {
		qp := f.QueryPath()
		if strings.HasPrefix(qp, fragment) {
			relativeToQ, _ := filepath.Rel(topCompleteDir, qp)
			remainingInPath := strings.Split(relativeToQ, string(os.PathSeparator))
			nextInPath := remainingInPath[0]
			fullCompletion := filepath.Join(topCompleteDir, nextInPath)

			if !resContains(fullCompletion) {
				compl := model.AutocompleteResult{
					Str:   fullCompletion,
					IsDir: len(remainingInPath) > 1,
				}
				res = append(res, compl)
			}
		}
	}

	if len(res) == 0 {
		return mc.suggestAutocompletionsLenient(fragment, allFiles)
		// didn't match anything so try a more flexible search
	}

	return res
}

func (mc *MainController) suggestAutocompletionsLenient(fragment string, allFiles []model.FilePath) []model.AutocompleteResult {
	res := make([]model.AutocompleteResult, 0)

	for _, f := range allFiles {
		qp := f.QueryPath()
		fileName := filepath.Base(qp)
		if strings.HasPrefix(fileName, fragment) {
			compl := model.AutocompleteResult{
				Str:   qp,
				IsDir: false,
			}
			res = append(res, compl)
		}
	}

	return res
}

func (mc *MainController) updateInput() {
	matched, res := mc.CompletionView.Current()
	if matched {
		str := res.Str
		if res.IsDir {
			str += "/"
		}
		mc.InputView.SetTextContentString(str)
		mc.InputView.SetCursorX(runewidth.StringWidth(str))
	}
}

func (mc *MainController) handleCompltionModeEvent(e *egg.KeyEvent) {
	// ensure conditions are correct
	if mc.activeMode == constants.ActiveModeAutocomplete && mc.CompletionView.IsOpen() {
		e.SetPropagate(false)

		switch e.Key {
		case egg.KeyTab, egg.KeyDown:
			mc.CompletionView.Next()
			mc.updateInput()
		case egg.KeyUp, egg.KeyBacktab:
			mc.CompletionView.Prev()
			mc.updateInput()
		case egg.KeyEnter:
			mc.setMode(constants.ActiveModeDefault)
			// nothing else
		default:
			mc.setMode(constants.ActiveModeDefault)
			e.SetPropagate(true)
			mc.handleEventInputMode(e)
		}
	}
}
