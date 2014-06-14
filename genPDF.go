package main

import (
	"fmt"
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
	doc := pdf.New()
	canvas := doc.NewPage(pdf.USLetterWidth, pdf.USLetterHeight)
	labels := avery22805.Positions()
	for i, label := range labels {
		qrCode, err := qr.Encode(fmt.Sprintf("http://foo/bar/%d", i), qr.H)
		if err != nil {
			panic(err)
		}
		image := qrCode.Image()
		canvas.DrawImage(image, label)
	}
	canvas.Close()

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
	y := (l.PageHeight / two) - pdf.Unit(l.Rows / 2) * (l.Height + l.RowGap) + l.RowGap / two
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
		y += l.Height + l.RowGap
	}
	return labels
}
