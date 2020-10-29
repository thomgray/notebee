package view

import (
	"log"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/thomgray/egg"
	"github.com/thomgray/notebee/model"
)

type CompletionView struct {
	*egg.View
	completions   []model.AutocompleteResult
	selectedIndex int
	open          bool
}

func MakeCompletionView() *CompletionView {
	cv := CompletionView{
		View:          egg.MakeView(),
		completions:   make([]model.AutocompleteResult, 0),
		selectedIndex: -1,
	}

	cv.SetVisible(false)
	cv.OnDraw(cv.draw)
	cv.Resize()

	egg.GetApplication().AddViewController(cv)

	return &cv
}

func (cv *CompletionView) Close() {
	cv.open = false
	cv.selectedIndex = -1
	cv.SetVisible(false)
}

func (cv *CompletionView) Open() {
	cv.open = true
	cv.selectedIndex = -1
	cv.SetVisible(true)
}

func (cv *CompletionView) IsOpen() bool {
	return cv.open
}

func (cv CompletionView) GetView() *egg.View {
	return cv.View
}

func (cv *CompletionView) SetCompletions(compl []model.AutocompleteResult) {
	cv.completions = compl
	cv.Resize()
}

func (cv *CompletionView) MaxHeight() int {
	return egg.WindowHeight() - 1
}

func (cv *CompletionView) Resize() {
	w := egg.WindowWidth()
	h := len(cv.completions) + 1

	log.Printf("completion view resized w=%d, h=%d", w, h)
	newBounds := egg.MakeBounds(0, 1, w, h)
	newBounds.Height = h
	newBounds.Width = w
	cv.SetBounds(newBounds)
}

func (cv *CompletionView) Next() {
	cv.selectedIndex++
	if cv.selectedIndex >= len(cv.completions) {
		cv.selectedIndex = len(cv.completions) - 1
	}
}

func (cv *CompletionView) Current() (match bool, res model.AutocompleteResult) {
	if cv.selectedIndex >= 0 && cv.selectedIndex < len(cv.completions) {
		res := cv.completions[cv.selectedIndex]
		return true, res
	}
	return false, model.AutocompleteResult{}
}

func (cv *CompletionView) Prev() {
	cv.selectedIndex--
	if cv.selectedIndex < 0 {
		cv.selectedIndex = 0
	}
}

func (cv *CompletionView) draw(c egg.Canvas) {
	selectedFg := egg.ColorBlack
	selectedBg := egg.ColorBlue
	log.Printf("drawing completion view, completions=%d", len(cv.completions))
	h := cv.MaxHeight()
	drawElipse := false
	if cv.completions != nil {
		for i, compl := range cv.completions {
			isSelected := i == cv.selectedIndex
			if i >= h-1 {
				drawElipse = true
				break
			}

			bg := c.Background
			if isSelected {
				bg = selectedBg
			}

			pieces := strings.Split(compl.Str, string(os.PathSeparator))
			lastPieceI := len(pieces) - 1
			x := 0
			for ii, piece := range pieces {
				fg := egg.ColorCyan
				isFinalPiece := ii == lastPieceI
				if isSelected {
					fg = selectedFg
				} else if isFinalPiece && !compl.IsDir {
					fg = c.Foreground
				}
				c.DrawString(piece, x, i, fg, bg, c.Attribute)

				x += runewidth.StringWidth(piece)
				if !isFinalPiece || compl.IsDir {
					slashFg := egg.ColorBrightMagenta
					if isSelected {
						slashFg = selectedFg
					}
					c.DrawRune('/', x, i, slashFg, bg, c.Attribute)
					x++
				}
			}
		}
	}
	if drawElipse {
		c.DrawString2("...", 0, h-1)
	}
	y := len(cv.completions)

	c.DrawString(strings.Repeat("─", c.Width), 0, y, egg.ColorBlue, c.Background, c.Attribute)

}
