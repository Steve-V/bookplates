package main

import (
	"bitbucket.org/zombiezen/gopdf/pdf"
	"code.google.com/p/rsc/qr"
	"fmt"
	"github.com/dichro/bookplates/format"
	"math"
	"os"
)

func main() {
	w := format.NewLabelSheetWriter(format.Avery22805)
	w.Write(&simple{590, 12729239}, 24)
	doc := w.Finish()
	if err := doc.Encode(os.Stdout); err != nil {
		panic(err)
	}
}

type simple struct {
	group, id int
}

func (s *simple) Render(page *pdf.Canvas, bound pdf.Point, index int) {
	qrCode, err := qr.Encode(fmt.Sprintf("http://bcing.me/%d-%d", s.group, s.id+index), qr.L)
	if err != nil {
		panic(err)
	}
	textSize := pdf.Unit(8)
	qrScale := float32(bound.X-(1.25*textSize)) / float32(bound.X)
	page.Push()
	page.Translate(1.25*textSize, 0)
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
	page.Rotate(math.Pi / 2)
	text = new(pdf.Text)
	text.SetFont(pdf.Helvetica, textSize)
	text.Text(fmt.Sprintf("BCID %d-%d", s.group, s.id+index))
	page.DrawText(text)
	page.Pop()
}
