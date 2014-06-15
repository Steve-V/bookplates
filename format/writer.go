package bookplates

import (
	"bitbucket.org/zombiezen/gopdf/pdf"
)

var Avery22805 = LabelSheet{
	Width: 1.5*pdf.Inch,
	Height: 1.5*pdf.Inch,
	PageWidth: 8.5*pdf.Inch,
	PageHeight: 11*pdf.Inch,
	Cols: 4,
	ColGap: 0.3125*pdf.Inch,
	Rows: 6,
	RowGap: 0.2*pdf.Inch,
}

type LabelSheet struct {
	// Label size.
	Width, Height pdf.Unit
	// Page size of sticker sheet.
	PageWidth, PageHeight pdf.Unit
	// Gap sizes between adjacent stickers.
	RowGap, ColGap pdf.Unit
	// Distribution of stickers on sheet.
	Rows, Cols int
}

// Assumes that label grid is centered in sheet and that the number of rows and columns is even.
func (l *LabelSheet) Positions() []pdf.Rectangle {
	labels := []pdf.Rectangle{}
	two := pdf.Unit(2)
	xBase := (l.PageWidth / two) - pdf.Unit(l.Cols / 2) * (l.Width + l.ColGap) + l.ColGap / two
	y := (l.PageHeight / two) + pdf.Unit(l.Rows / 2 - 1) * (l.Height + l.RowGap) + l.RowGap / two
	for row := 0; row < l.Rows; row++ {
		x := xBase
		for col := 0; col < l.Cols; col++ {
			labels = append(labels, pdf.Rectangle{
				Min: pdf.Point{
					X: x,
					Y: y,
				},
				Max: pdf.Point{
					X: x + l.Width,
					Y: y + l.Height,
				},
			})
			x += l.Width + l.ColGap
		}
		y -= l.Height + l.RowGap
	}
	return labels
}

type Renderer interface {
	Render(page *pdf.Canvas, bound pdf.Point, index int)
}

type LabelSheetWriter struct {
	paper LabelSheet
	positions []pdf.Rectangle
	doc *pdf.Document
	page *pdf.Canvas
	written int
}

func NewLabelSheetWriter(paper LabelSheet) *LabelSheetWriter {
	pos := paper.Positions()
	w := &LabelSheetWriter{
		paper: paper,
		positions: pos,
		doc: pdf.New(),
		written: len(pos),
	}
	return w
}

func (w *LabelSheetWriter) Write(renderer Renderer, count int) {
	for i := 0; i < count; i++ {
		w.write(renderer, i)
	}
}

func (w *LabelSheetWriter) write(renderer Renderer, index int) {
	if w.written == len(w.positions) {
		if w.page != nil {
			w.page.Close()
		}
		w.page = w.doc.NewPage(w.paper.PageWidth, w.paper.PageHeight)
		w.written = 0
	}
	pos := w.positions[w.written]
	w.page.Push()
	w.page.Translate(pos.Min.X, pos.Min.Y)
	w.written++
	renderer.Render(w.page, pdf.Point{
			X: w.paper.Width,
			Y: w.paper.Height,
		}, index)
	w.page.Pop()
}

func (w *LabelSheetWriter) Finish() *pdf.Document {
	if w.page != nil {
		w.page.Close()
		w.written = len(w.positions)
	}
	return w.doc
}

