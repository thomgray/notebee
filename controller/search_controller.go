package controller

import (
	"log"
	"strings"

	"github.com/thomgray/egg"
	"github.com/thomgray/notebee/constants"
	"github.com/thomgray/notebee/model"
)

func (mc *MainController) handleSearch(str string) {
	defer app.ReDraw()
	var scoredFiles []*model.SearchResultItem
	allPaths := mc.FileManager.FindSupportedFilePaths()
	for _, p := range allPaths {
		f := model.LoadCodeFile(p.Full)
		content := string(f.Content)
		if strings.Contains(content, str) {
			scoredFiles = append(scoredFiles, &model.SearchResultItem{
				File:  f,
				Path:  p,
				Score: 1,
			})
		}
	}

	for _, s := range scoredFiles {
		log.Printf("scored file name=%s path=%s", s.Path.Relative, s.Path.QueryPath())
	}

	if len(scoredFiles) > 0 {
		mc.setMode(constants.ActiveModeSearchResultSelect)
		mc.SearchResultsView.SetItems(scoredFiles)
		mc.SearchResultsView.Open()
	}
}

func (mc *MainController) handleSearchResultModeEvent(e *egg.KeyEvent) {
	e.SetPropagate(false)
	switch e.Key {
	case egg.KeyUp, egg.KeyBacktab:
		mc.SearchResultsView.Prev()
	case egg.KeyDown, egg.KeyTAB:
		mc.SearchResultsView.Next()
	case egg.KeyEnter:
		result := mc.SearchResultsView.Selected()
		if result != nil {
			mc.setMode(constants.ActiveModeDefault)
			mc.setInputMode(constants.InputModeTraverse)
			mc.handleTraverse(result.Path.QueryPath())
		}
	default:
		mc.setMode(constants.ActiveModeDefault)
		e.SetPropagate(true)
		mc.handleEventInputMode(e)
	}
}
