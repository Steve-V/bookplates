package main

import (
	"fmt"
	"os"
	"code.google.com/p/rsc/qr"
	"bitbucket.org/zombiezen/gopdf/pdf"
)

func main() {
	doc := pdf.New()
	canvas := doc.NewPage(pdf.USLetterWidth, pdf.USLetterHeight)
	for row := 0*pdf.Inch; row <= pdf.USLetterHeight-1; row += pdf.Inch {
		for col := 0*pdf.Inch; col <= pdf.USLetterWidth-1; col += pdf.Inch {
			qrCode, err := qr.Encode(fmt.Sprintf("http://foo/bar/%d/%d", row, col), qr.H)
			if err != nil {
				panic(err)
			}

			canvas.DrawImage(qrCode.Image(), pdf.Rectangle{
				Min: pdf.Point{X: col, Y: row},
				Max: pdf.Point{X: col+pdf.Inch, Y: row+pdf.Inch},
			})
		}
	}
	canvas.Close()

	if err := doc.Encode(os.Stdout); err != nil {
		panic(err)
	}
}
