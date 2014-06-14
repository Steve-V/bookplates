package main

import (
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/signintech/gopdf/fonts"
)

func main() {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddFont("THSarabunPSK", new(fonts.THSarabun), "THSarabun.z")
	pdf.AddFont("Loma", new(fonts.Loma), "Loma.z")
	pdf.AddPage()
	pdf.SetFont("THSarabunPSK", "B", 14)
	pdf.Cell(nil, ToCp874("Hello world  = สวัสดี โลก in thai"))
	pdf.WritePdf("/tmp/out.pdf")
	fmt.Println("Done...")
}

func ToCp874(str string) string {
	return str
}
