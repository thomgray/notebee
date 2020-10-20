package view

import (
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/thomgray/egg"
	"github.com/thomgray/notebee/model"
)

type CompletionView struct {
	*egg.View
	completions []model.AutocompleteResult
}

func MakeCompletionView() *CompletionView {
	cv := CompletionView{
		View:        egg.MakeView(),
		completions: make([]model.AutocompleteResult, 0),
	}

	cv.SetVisible(false)
	cv.OnDraw(cv.draw)
	cv.Resize()

	egg.GetApplication().AddViewController(cv)

	return &cv
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

	newBounds := egg.MakeBounds(0, 1, w, h)
	newBounds.Height = h
	newBounds.Width = w
	cv.SetBounds(newBounds)
}

func (cv *CompletionView) draw(c egg.Canvas, _ egg.State) {
	h := cv.MaxHeight()
	drawElipse := false
	if cv.completions != nil {
		for i, compl := range cv.completions {
			if i >= h-1 {
				drawElipse = true
				break
			}

			pieces := strings.Split(compl.Str, string(os.PathSeparator))
			lastPieceI := len(pieces) - 1
			x := 0
			for ii, piece := range pieces {
				fg := egg.ColorCyan
				isFinalPiece := ii == lastPieceI
				if isFinalPiece && !compl.IsDir {
					fg = c.Foreground
				}
				c.DrawString(piece, x, i, fg, c.Background, c.Attribute)

				x += runewidth.StringWidth(piece)
				if !isFinalPiece || compl.IsDir {
					c.DrawRune('/', x, i, egg.ColorMagenta, c.Background, c.Attribute)
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
