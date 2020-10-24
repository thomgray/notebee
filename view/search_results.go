package view

import (
	"log"

	"github.com/thomgray/egg"
	"github.com/thomgray/notebee/model"
)

type SearchResultsView struct {
	*egg.View
	Items     []*model.SearchResultItem
	itemIndex int
	open      bool
}

func (sv SearchResultsView) New() *SearchResultsView {
	sv.View = egg.MakeView()
	sv.OnKeyEvent(sv.handleKeyEvent)
	sv.OnDraw(sv.draw)
	sv.itemIndex = -1
	return &sv
}

func (sv *SearchResultsView) Refit(w, h int) {
	anchor := 1
	sv.SetBounds(egg.MakeBounds(0, anchor, w, h-anchor))
}

func (sv *SearchResultsView) handleKeyEvent(e *egg.KeyEvent) {
	log.Println("Key event biatch")
}

func (sv *SearchResultsView) SetItems(items []*model.SearchResultItem) {
	bnds := sv.GetBounds()
	bnds.Height = len(items)
	sv.SetBounds(bnds)
	sv.Items = items
}

func (sv *SearchResultsView) Open() {
	sv.open = true
	sv.SetVisible(sv.open)
}

func (sv *SearchResultsView) Close() {
	sv.open = false
	sv.SetVisible(sv.open)
	sv.itemIndex = -1
	sv.Items = []*model.SearchResultItem{}
}

func (sv *SearchResultsView) Selected() *model.SearchResultItem {
	if sv.itemIndex >= 0 && sv.itemIndex < len(sv.Items) {
		return sv.Items[sv.itemIndex]
	}
	return nil
}

func (sv *SearchResultsView) IsOpen() bool {
	return sv.open
}

func (sv *SearchResultsView) Next() {
	sv.itemIndex++
	if sv.itemIndex >= len(sv.Items) {
		sv.itemIndex = len(sv.Items) - 1
	}
}

func (sv *SearchResultsView) Prev() {
	sv.itemIndex--
	if sv.itemIndex < 0 {
		sv.itemIndex = 0
	}
}

func (sv *SearchResultsView) draw(c egg.Canvas) {
	for i, res := range sv.Items {
		bg := c.Background
		if i == sv.itemIndex {
			bg = egg.ColorBrightCyan
		}
		c.DrawString(res.Path.QueryPath(), 0, i, c.Foreground, bg, c.Attribute)
	}
}
