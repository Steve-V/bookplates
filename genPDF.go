package main

import (
	"fmt"
	"math"
	"os"
	"code.google.com/p/rsc/qr"
	"bitbucket.org/zombiezen/gopdf/pdf"
)

var avery22805 = labelSheet{
	Width: 1.5*pdf.Inch,
	Height: 1.5*pdf.Inch,
	PageWidth: 8.5*pdf.Inch,
	PageHeight: 11*pdf.Inch,
	Cols: 4,
	ColGap: 0.3125*pdf.Inch,
	Rows: 6,
	RowGap: 0.2*pdf.Inch,
}

func main() {
	w := newLabelSheetWriter(avery22805)
	group, id := 590, 12729239
	for i := 0; i < 24; i++ {
		w.Write(group, id + i)
	}
	doc := w.Finish()
	if err := doc.Encode(os.Stdout); err != nil {
		panic(err)
	}
}

type labelSheet struct {
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
func (l *labelSheet) Positions() []pdf.Rectangle {
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

type labelSheetWriter struct {
	paper labelSheet
	positions []pdf.Rectangle
	doc *pdf.Document
	page *pdf.Canvas
	written int
}

func newLabelSheetWriter(paper labelSheet) *labelSheetWriter {
	pos := paper.Positions()
	w := &labelSheetWriter{
		paper: paper,
		positions: pos,
		doc: pdf.New(),
		written: len(pos),
	}
	return w
}

func (w *labelSheetWriter) Write(group, id int) {
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
	render(w.page, pdf.Point{
			X: w.paper.Width,
			Y: w.paper.Height,
		}, group, id)
	w.page.Pop()
}

func render(page *pdf.Canvas, bound pdf.Point, group, id int) {
	qrCode, err := qr.Encode(fmt.Sprintf("http://bcing.me/%d-%d", group, id), qr.L)
	if err != nil {
		panic(err)
	}
	textSize := pdf.Unit(8)
	qrScale := float32(bound.X - (1.25 * textSize)) / float32(bound.X)
	page.Push()
	page.Translate(1.25 * textSize, 0)
	page.Scale(qrScale, qrScale)
	page.DrawImage(qrCode.Image(), pdf.Rectangle{
		Max: bound,
	})
	page.Pop()
	page.Push()
	text := new(pdf.Text)
	text.SetFont(pdf.Helvetica, textSize)
	text.Text("Ex libris Miki dichro@rcpt.to")
	page.Translate(bound.X-text.X(), bound.Y-textSize)
	page.DrawText(text)
	page.Pop()
	page.Push()
	page.Translate(textSize, 0)
	page.Rotate(math.Pi/2)
	text = new(pdf.Text)
	text.SetFont(pdf.Helvetica, textSize)
	text.Text(fmt.Sprintf("BCID %d-%d", group, id))
	page.DrawText(text)
	page.Pop()
}

func (w *labelSheetWriter) Finish() *pdf.Document {
	if w.page != nil {
		w.page.Close()
		w.written = len(w.positions)
	}
	return w.doc
}
