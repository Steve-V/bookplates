package main

import (
	"os"
	"code.google.com/p/rsc/qr"
	"bitbucket.org/zombiezen/gopdf/pdf"
)

func main() {
	qrCode, err := qr.Encode("http://foo/bar", qr.H)
	if err != nil {
		panic(err)
	}

	doc := pdf.New()
	canvas := doc.NewPage(pdf.USLetterWidth, pdf.USLetterHeight)
	canvas.DrawImage(qrCode.Image(), pdf.Rectangle{
		Min: pdf.Point{X: 7.5*pdf.Inch, Y: 7.5*pdf.Inch},
		Max: pdf.Point{X: 8.5*pdf.Inch, Y: 8.5*pdf.Inch},
	})
	canvas.Close()

	if err := doc.Encode(os.Stdout); err != nil {
		panic(err)
	}
}
