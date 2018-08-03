package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

const (
	bannerHt = 95.0
	xIndent  = 40.0
)

func main() {
	pdf := gofpdf.New(gofpdf.OrientationLandscape, gofpdf.UnitPoint, gofpdf.PageSizeLetter, "")
	pdf.AddPage()

	//pdf = drawBasic(pdf)
	pdf = drawBanner(pdf)

	_, lineHeight := pdf.GetFontSize()
	drawSummaryBlock(pdf, xIndent, bannerHt+lineHeight*2.0, "Billed To", "Client Name", "123 Client", "State, Country", "Postal Code")
	drawSummaryBlock(pdf, xIndent*2.0+lineHeight*11.0, bannerHt+lineHeight*2.0, "Invoice Number", "0000000000012")

	//drawGrid
	drawGrid(pdf)
	err := pdf.OutputFileAndClose("p1.pdf")
	if err != nil {
		panic(err)
	}
}

func getLongText() string {
	text := "Here is some text. If it is too long it will be word wrapped automatically. If there is a new line it will be\nwrapped as well (unlike other ways of writing text in gofpdf)."
	return text
}

func drawGrid(pdf *gofpdf.Fpdf) {
	w, h := pdf.GetPageSize()
	pdf.SetFont("courier", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetDrawColor(200, 200, 200)
	for x := 0.0; x < w; x = x + (w / 20.0) {
		pdf.SetTextColor(200, 200, 200)
		pdf.Line(x, 0, x, h)
		_, lineHeight := pdf.GetFontSize()
		pdf.Text(x, lineHeight, fmt.Sprintf("%d", int(x)))
	}

	for y := 0.0; y < h; y = y + (w / 20.0) {
		if y < bannerHt*0.9 {
			pdf.SetTextColor(200, 200, 200)
		} else {
			pdf.SetTextColor(80, 80, 80)
		}
		pdf.Line(0, y, w, y)
		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
	}
}

func drawBasic(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	//basic text stuff
	pdf.SetFont("Arial", "B", 28)
	_, lineHeight := pdf.GetFontSize()
	pdf.SetTextColor(255, 0, 0)
	pdf.Text(0, lineHeight, "Hello, world")
	pdf.MoveTo(0, lineHeight*2)

	pdf.SetFont("times", "", 18)
	pdf.SetTextColor(100, 100, 100)
	_, lineHeight = pdf.GetFontSize()
	pdf.MultiCell(0, lineHeight*1.5, getLongText(), gofpdf.BorderNone, gofpdf.AlignMiddle, false)

	//basic shape stuff
	pdf.SetFillColor(0, 255, 0)
	pdf.SetDrawColor(0, 0, 255)
	pdf.Rect(10, 150, 100, 100, "FD")
	pdf.SetFillColor(100, 200, 200)
	pdf.Polygon([]gofpdf.PointType{
		{110, 250},
		{160, 300},
		{110, 350},
		{60, 300},
	}, "F")

	//image
	pdf.ImageOptions("images/logo.png", 275, 275, 192, 32, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")
	return pdf
}

func drawBanner(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	w, h := pdf.GetPageSize() //get relative value
	fmt.Printf("width=%v, height=%v\n", w, h)
	pdf.SetFillColor(103, 60, 79) //dark maroon
	//top banner
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{w, 0},
		{w, bannerHt},
		{0, bannerHt * 0.9},
	}, "F")
	//bottom banner
	pdf.Polygon([]gofpdf.PointType{
		{0, h},
		{0, h - bannerHt*0.2},
		{w, h - bannerHt*0.1},
		{w, h},
	}, "F")
	//Banner - invoice
	pdf.SetFont("arial", "B", 40)
	pdf.SetTextColor(255, 255, 255) //white
	_, lineHeight := pdf.GetFontSize()
	//whatever the lineHeit is, subtract the bannerHt divided by 2
	pdf.Text(xIndent, bannerHt/2.0+lineHeight/3.1, "INVOICE")

	//Banner - phone, e-mail, domain
	pdf.SetFont("arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-2.0*124.0, (bannerHt-lineHeight*1.5*3.0)/2.0)
	pdf.MultiCell(124.0, lineHeight*1.5, "(123) 456-78\njon@calhoun.io\nGophercises.com", gofpdf.BorderNone, gofpdf.AlignMiddle, false)

	//Banner - address and Co. info
	pdf.SetFont("arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-124.0, (bannerHt-lineHeight*1.5*3.0)/2.0)
	pdf.MultiCell(124.0, lineHeight*1.5, "Friedrich Strasse 68\n10117\nBerlin", gofpdf.BorderNone, gofpdf.AlignRight, false)
	return pdf
}

func drawSummaryBlock(pdf *gofpdf.Fpdf, x, y float64, title string, data ...string) (float64, float64) {
	//Content - billed to
	pdf.SetFont("times", "", 14)
	pdf.SetTextColor(180, 180, 180)
	_, lineHeight := pdf.GetFontSize()
	y = y + lineHeight
	pdf.Text(x, y, title)

	pdf.SetTextColor(50, 50, 50)
	y = y + lineHeight*.25
	for _, str := range data {
		y = y + lineHeight*1.25
		pdf.Text(x, y, str)
	}
	return x, y
}
